package runtimes

import (
	"github.com/mee6aas/kyle/internal/pkg/runtime"
)

var (
	//				   actName
	runtimes = make(map[string][]*runtime.Runtime)
)
