package agent

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mee6aas/zeep/api"
	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"

	"github.com/mee6aas/kyle/internal/pkg/cact"
	"github.com/mee6aas/kyle/internal/pkg/client/invokee"
	"github.com/mee6aas/kyle/internal/pkg/client/invoker"
	"github.com/mee6aas/kyle/internal/pkg/pool"
	assigns "github.com/mee6aas/kyle/internal/pkg/var/assignments"
)

// Start starts the agent
func Start(ctx context.Context) (e error) {
	var (
		actMain string // name of the activity allocated to this container
	)

	// wg := sync.WaitGroup{}
	// defer wg.Wait()

	startCtx, _ := context.WithCancel(ctx)
	// startCtx, startCancel := context.WithCancel(ctx)
	// defer startCancel()

	defer func() {
		if e == nil {
			return
		}

		log.WithError(e).Error("Agent failed")
	}()

	e = invokee.Connect(startCtx, localAgntAddr)
	if e != nil {
		e = errors.Wrapf(e, "Failed to connect to local agent invokee service at %s", localAgntAddr)
		return
	}

	e = invoker.Connect(startCtx, localAgntAddr)
	if e != nil {
		e = errors.Wrapf(e, "Failed to connect to local agent invoker service at %s", localAgntAddr)
		return
	}

	// wg.Add(1)
	go func() {
		// defer wg.Done()
		// defer startCancel()

		log.WithField("addr", serveAddr).Debug("Serving the agent")
		err := serve(startCtx)
		log.WithError(err).Error("Agent serve failed")
		if err != nil && err != context.Canceled {
			e = errors.Wrap(err, "Agent serve failed")
			log.WithError(e).Error("Agent serve failed")
			return
		}
	}()

	// wg.Add(1)
	go func() {
		// defer wg.Done()
		// defer startCancel()

		log.Debug("Starting the runtime pool")
		err := pool.Start(startCtx)
		log.WithError(err).Error("Pool failed")
		if err != nil && err != context.Canceled {
			e = errors.Wrap(err, "Pool failed")
			log.WithError(e).Error("Pool failed")
			return
		}
	}()

	log.Debug("Waiting the task::LOAD")
	t, e := invokee.FetchTask(startCtx)
	if e != nil {
		e = errors.Wrap(e, "Failed to fetch a task from the invokee service")
		return
	}

	if t.GetType() != invokeeV1API.TaskType_LOAD {
		e = errors.Wrap(e, "First task must be LOAD")
		return
	}

	actMain = t.GetArg()
	ap := filepath.Join(api.ActivityResource, actMain, api.ActivityManifestName)
	log.WithField("path", ap).Debug("Unmarshalling the activity manifest")
	e = cact.UnmarshalFromFile(ap)
	if e != nil {
		e = errors.Wrapf(e, "Failed to unmarshal the activity manifest at %s", ap)
		return
	}

	log.Debug("Fetch the runtime")
	r, e := pool.Fetch(startCtx)
	if e != nil {
		e = errors.Wrap(e, "Failed to fetch the runtime")
		return
	}
	pid, ok := r.PID()
	if !ok {
		e = errors.New("Failed to get PID from the runtime")
		return
	}

	id, onResolved := assigns.Add()

	log.WithFields(log.Fields{
		"arg": actMain,
		"pid": pid,
	}).Debug("Assign the task LOAD to the runtime")
	e = r.Assign(startCtx, invokeeV1API.Task{
		Id:   id,
		Type: invokeeV1API.TaskType_LOAD,
		Arg:  actMain,
	})
	if e != nil {
		e = errors.Wrap(e, "Failed to assign a task to the runtime")
		return
	}

	isWorkflow := cact.HasDep()

	log.WithField("workflow", isWorkflow).Debug("Is the task a workflow?")

	if !isWorkflow {
		log.Debug("Starting handover procedure")

		// notify to the local agent
		e = invokee.Handover(startCtx)
		if e != nil {
			e = errors.Wrap(e, "Failed to request handover")
			return
		}

		// handover
		log.WithField("arg", localAgntAddr).Debug("Assign task::HANDOVER to the runtime")
		r.Assign(startCtx, invokeeV1API.Task{
			Type: invokeeV1API.TaskType_HANDOVER,
			Arg:  localAgntAddr,
		})

		_, e = invokee.FetchTask(startCtx)
		if e == nil {
			e = errors.Wrap(e, "Expected that the connection from the local agent is canceled")
			return
		}
		if e == context.Canceled {
			return
		}

		log.Info("The control is handovered to the local agent")
	}

	// forward report for LOAD task
	{
		resolved := <-onResolved
		log.Debug("task::Load resolved")
		rst := resolved.(*invokeeV1API.ReportRequest)

		e = invokee.Report(startCtx, t.GetId(), rst.GetResult())
		if e != nil {
			e = errors.Wrap(e, "Failed to report a task::LOAD to invokee service")
			return
		}
	}

	if !isWorkflow {
		log.Debug("Stopping gRPC server")
		gRPCServer.GracefulStop()

		return
	}

	log.Info("Listen tasks")
	for {
		log.Info("Wait for the task")
		t, e := invokee.FetchTask(startCtx)
		if e != nil {
			e = errors.Wrap(e, "Failed to fetch a task from the invokee service")
			return e
		}
		if ttype := t.GetType(); ttype != invokeeV1API.TaskType_INVOKE {
			e = errors.Wrapf(e, "Expected that the type of the task is INVOKE but %s", ttype.String())
			return e
		}

		id, onInvoked := assigns.Add()
		e = r.Assign(startCtx, invokeeV1API.Task{
			Id:   id,
			Type: invokeeV1API.TaskType_INVOKE,
			Arg:  t.GetArg(),
		})

		// res, e := ivkerV1Hndl.InvokeRequested(startCtx, nil, "", actMain, t.GetArg())
		if e != nil {
			e = errors.Wrap(e, "Failed to invoke")
			return e
		}

		res, ok := <-onInvoked
		if !ok {
			panic("Runtime disconnected while invoke an activity")
		}

		// FIXME
		// resolve task version
		rst := res.(*invokeeV1API.ReportRequest).GetResult()

		e = invokee.Report(startCtx, t.GetId(), rst)
		if e != nil {
			e = errors.Wrap(e, "Failed to report")
			return e
		}
	}
}
