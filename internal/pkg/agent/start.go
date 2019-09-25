package agent

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mee6aas/zeep/api"
	"github.com/mee6aas/zeep/pkg/activity"
	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"

	"github.com/mee6aas/kyle/internal/pkg/client/invokee"
	"github.com/mee6aas/kyle/internal/pkg/pool"
	assigns "github.com/mee6aas/kyle/internal/pkg/var/assignments"
)

// Start starts the agent
func Start(ctx context.Context) (e error) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	startCtx, startCancel := context.WithCancel(ctx)
	defer startCancel()

	defer func() {
		if e == nil {
			return
		}

		log.WithError(e).Error("agent.start failed")
	}()

	e = invokee.Connect(startCtx, localAgntAddr)
	if e != nil {
		e = errors.Wrapf(e, "Failed to connect to local agent invokee service at %s", localAgntAddr)
		return
	}

	wg.Add(1)
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

	wg.Add(1)
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

	log.Debug("Waiting the task LOAD")
	t, e := invokee.FetchTask(startCtx)
	if e != nil {
		e = errors.Wrap(e, "Failed to fetch a task from the invokee service")
		return
	}

	if t.GetType() != invokeeV1API.TaskType_LOAD {
		e = errors.Wrap(e, "First task must be LOAD")
		return
	}

	ap := filepath.Join(api.ActivityResource, t.GetArg(), api.ActivityManifestName)
	log.WithField("path", ap).Debug("Unmarshalling the activity manifest")
	a, e := activity.UnmarshalFromFile(ap)
	if e != nil {
		e = errors.Wrapf(e, "Failed to unmarshal the activity manifest at %s", ap)
		return
	}

	log.Debug("Fetch the runtime")
	r, e := pool.Fetch(startCtx)
	if e != nil {
		e = errors.Wrap(e, "Failed to fetch runtime")
		return
	}
	pid, ok := r.PID()
	if !ok {
		panic("Failed to get ip form the runtime")
	}

	id, onResolved := assigns.Add()

	log.WithFields(log.Fields{
		"arg": t.GetArg(),
		"pid": pid,
	}).Debug("Give task LOAD to the runtime")
	e = r.Assign(startCtx, invokeeV1API.Task{
		Id:   id,
		Type: invokeeV1API.TaskType_LOAD,
		Arg:  t.GetArg(),
	})
	if e != nil {
		e = errors.Wrap(e, "Failed to assign a task to the runtime")
		return
	}

	isWorkflow := len(a.Dependencies) > 0

	log.WithField("workflow", isWorkflow).Debug("Is the task a workflow?")

	if !isWorkflow {
		log.Debug("Starting handover procedure")

		// notify to local agent
		e = invokee.Handover(startCtx)
		if e != nil {
			e = errors.Wrap(e, "Failed request to handover")
			return
		}

		// handover
		log.WithField("arg", localAgntAddr).Debug("Give task HANDOVER to the runtime")
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

		// handovered
	}

	// forward report for LOAD task
	{
		resolved := <-onResolved
		log.Debug("resolved?")
		rst := resolved.(*invokeeV1API.ReportRequest)

		e = invokee.Report(startCtx, t.GetId(), rst.GetResult())
		if e != nil {
			e = errors.Wrap(e, "Failed to report a LOAD task to invokee service")
			return
		}
	}

	if !isWorkflow {
		gRPCServer.GracefulStop()

		return
	}

	wg.Wait()

	return
}
