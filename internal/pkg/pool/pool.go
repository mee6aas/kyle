package pool

import (
	"context"

	"github.com/mee6aas/kyle/internal/pkg/runtime"
)

var (
	mngCtx    context.Context
	mngCancel context.CancelFunc

	rConf runtime.Config

	onFetched = make(chan struct{}, 1)
)
