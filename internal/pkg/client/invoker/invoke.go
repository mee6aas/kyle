package invoker

import (
	"context"

	v1 "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

// Invoke request to invoke the specified activity name with the arg.
func Invoke(ctx context.Context, actName string, arg string) (rst string, e error) {
	res, e := client.Invoke(ctx, &v1.InvokeRequest{
		ActName: actName,
		Arg:     arg,
	})

	if e != nil {
		return
	}

	rst = res.GetResult()

	return
}
