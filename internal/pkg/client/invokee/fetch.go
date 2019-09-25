package invokee

import (
	"context"

	"github.com/pkg/errors"

	api "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

// FetchTask withdraws task from task stream.
// it blocked until channel returns.
func FetchTask(ctx context.Context) (t api.Task, e error) {
	var (
		ok bool
	)

	select {
	case t, ok = <-onTask:
		if !ok {
			e = errors.New("Channel closed")
			return
		}

	case <-ctx.Done():
		e = ctx.Err()
		return
	}

	return
}
