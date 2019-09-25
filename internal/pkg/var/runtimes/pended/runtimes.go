package runtimes

import (
	"github.com/mee6aas/kyle/internal/pkg/runtime"
)

type elem struct {
	runtime   *runtime.Runtime
	onRelease chan struct{}
}

var (
	//			 		PID
	runtimes = make(map[int]elem)
)
