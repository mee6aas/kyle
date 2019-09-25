package invokee

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	v1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

// Connect connects to the agent with the specified agent.
func Connect(ctx context.Context, addr string) (e error) {
	if conn, e = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock()); e != nil {
		e = errors.Wrap(e, "Failed to dial")
		return
	}

	client = v1.NewInvokeeClient(conn)

	if e = listen(ctx); e != nil {
		e = errors.Wrap(e, "Failed to listen")
		return
	}

	return
}
