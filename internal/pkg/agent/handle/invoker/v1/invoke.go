package v1

import (
	"context"
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	invokerV1API "github.com/mee6aas/zeep/pkg/api/invoker/v1"
	invokerV1 "github.com/mee6aas/zeep/pkg/service/invoker/v1"

	"github.com/mee6aas/kyle/internal/pkg/cact"
	"github.com/mee6aas/kyle/internal/pkg/client/invoker"
	"github.com/mee6aas/kyle/internal/pkg/pool"
	assigns "github.com/mee6aas/kyle/internal/pkg/var/assignments"
	rsAllocated "github.com/mee6aas/kyle/internal/pkg/var/runtimes/allocated"
)

// InvokeRequested is invoked when the invoker requests an activity invoke.
func (h Handle) InvokeRequested(
	ctx context.Context,
	_ *net.TCPAddr,
	_ string,
	actName string,
	arg string,
) (
	res *invokerV1.InvokeResponse,
	e error,
) {
	defer func() {
		if e == nil {
			return
		}

		log.WithFields(log.Fields{
			"actName": actName,
			"arg":     arg,
		}).Warn("Invoker client returns error")
	}()

	if actName == "" {
		e = status.Error(codes.InvalidArgument, "Name of the activity to invoke not provided")
		return
	}

	dep, ok := cact.Dep(actName)
	if !ok {
		e = status.Error(codes.NotFound, "Activity not exists in the dependency list")
		return
	}

	if dep.Outflow == "always" {
		rst, e := invoker.Invoke(ctx, actName, arg)
		return &invokerV1API.InvokeResponse{Result: rst}, e
	}

	r, e := rsAllocated.Take(actName)
	if e != nil {
		r, e = pool.Fetch(ctx)
		if e != nil {
			panic(e)
		}

		id, onLoaded := assigns.Add()

		if e = r.Assign(ctx, invokeeV1API.Task{
			Id:   id,
			Type: invokeeV1API.TaskType_LOAD,
			Arg:  actName,
		}); e != nil {
			e = errors.Wrap(e, "Failed to assign task::LOAD to runtime")
			return
		}

		// TODO: use select with ctx
		_, ok := <-onLoaded
		if !ok {
			panic("Runtime disconnected while load an activity")
		}
	}

	id, onInvoked := assigns.Add()

	if e = r.Assign(ctx, invokeeV1API.Task{
		Id:   id,
		Type: invokeeV1API.TaskType_INVOKE,
		Arg:  arg,
	}); e != nil {
		e = errors.Wrap(e, "Failed to assign task::INVOKE to runtime")
		return
	}

	// TODO: use select with ctx
	rst, ok := <-onInvoked
	if !ok {
		panic("Runtime disconnected while invoke an activity")
	}

	switch r := rst.(type) {
	case *invokeeV1API.ReportRequest:
		res = &invokerV1.InvokeResponse{
			Result: r.GetResult(),
		}
	default:
		panic("Unrecognized report request")
	}

	e = rsAllocated.Add(actName, r)

	return
}
