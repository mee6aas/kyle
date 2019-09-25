package runtimes

import (
	"github.com/mee6aas/kyle/internal/pkg/runtime"
)

var (
	//					PID
	runtimes = make(map[int]*runtime.Runtime)
)
