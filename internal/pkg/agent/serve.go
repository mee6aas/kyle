package agent

import (
	"context"

	"google.golang.org/grpc"

	server "github.com/mee6aas/zeep/pkg/protocol/grpc"

	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	invokerV1API "github.com/mee6aas/zeep/pkg/api/invoker/v1"
	invokeeV1Svc "github.com/mee6aas/zeep/pkg/service/invokee/v1"
	invokerV1Svc "github.com/mee6aas/zeep/pkg/service/invoker/v1"

	invokeeV1Handle "github.com/mee6aas/kyle/internal/pkg/agent/handle/invokee/v1"
	invokerV1Handle "github.com/mee6aas/kyle/internal/pkg/agent/handle/invoker/v1"
)

var (
	gRPCServer  *grpc.Server
	ivkeeV1Hndl = invokeeV1Handle.Handle{}
	ivkerV1Hndl = invokerV1Handle.Handle{}
)

func serve(ctx context.Context) (e error) {
	gRPCServer = grpc.NewServer()
	invokeeV1API.RegisterInvokeeServer(gRPCServer, invokeeV1Svc.NewInvokeeAPIServer(ivkeeV1Hndl))
	invokerV1API.RegisterInvokerServer(gRPCServer, invokerV1Svc.NewInvokerAPIServer(ivkerV1Hndl))

	if e = server.Serve(ctx, gRPCServer, serveAddr); e != nil {
		return
	}

	return
}
