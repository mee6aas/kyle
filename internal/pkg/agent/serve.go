package agent

import (
	"context"

	"google.golang.org/grpc"

	server "github.com/mee6aas/zeep/pkg/protocol/grpc"

	invokeeV1API "github.com/mee6aas/zeep/pkg/api/invokee/v1"
	invokeeV1Svc "github.com/mee6aas/zeep/pkg/service/invokee/v1"

	invokeeV1Handle "github.com/mee6aas/kyle/internal/pkg/agent/handle/invokee/v1"
)

var (
	gRPCServer *grpc.Server
)

func serve(ctx context.Context) (e error) {
	gRPCServer = grpc.NewServer()
	invokeeV1API.RegisterInvokeeServer(gRPCServer, invokeeV1Svc.NewInvokeeAPIServer(
		invokeeV1Handle.Handle{},
	))

	if e = server.Serve(ctx, gRPCServer, serveAddr); e != nil {
		return
	}

	return
}
