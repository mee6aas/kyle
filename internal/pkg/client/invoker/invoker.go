package invoker

import (
	"google.golang.org/grpc"

	api "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

var (
	client api.InvokerClient
	conn   *grpc.ClientConn
)
