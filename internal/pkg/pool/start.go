package pool

import (
	"context"
)

// Start start the management for the pool.
func Start(ctx context.Context) (e error) {
	mngCtx, mngCancel = context.WithCancel(ctx)
	defer mngCancel()

	e = spawn(mngCtx, rConf)
	if e != nil {
		return
	}

	for {
		select {
		case <-mngCtx.Done():
			e = ctx.Err()
			return
		case <-onFetched:
			e = spawn(mngCtx, rConf)
		}

		if e != nil {
			return
		}
	}
}
