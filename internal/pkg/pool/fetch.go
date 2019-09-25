package pool

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mee6aas/kyle/internal/pkg/runtime"
	runtimes "github.com/mee6aas/kyle/internal/pkg/var/runtimes/connected"
)

// Fetch withdraws a runtime in the pool.
func Fetch(ctx context.Context) (r *runtime.Runtime, e error) {
	defer func() {
		if e == nil {
			return
		}
		// go func() {
		// 	onFetched <- struct{}{}
		// }()
	}()

	r, ok := runtimes.Take()
	if ok {
		return
	}

	if e = spawn(ctx, rConf); e != nil {
		e = errors.Wrap(e, "Failed to spawn a runtime")
		return
	}

	return
}
