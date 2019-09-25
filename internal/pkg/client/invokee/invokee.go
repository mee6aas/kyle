package invokee

import (
	"google.golang.org/grpc"

	api "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

var (
	client api.InvokeeClient
	conn   *grpc.ClientConn
)
