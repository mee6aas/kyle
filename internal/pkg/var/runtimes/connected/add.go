package runtimes

import (
	"github.com/mee6aas/kyle/internal/pkg/runtime"
)

// Add adds a runtime into the collection.
func Add(r *runtime.Runtime) bool {
	pid, ok := r.PID()
	if !ok {
		return false
	}

	if _, ok = runtimes[pid]; ok {
		return false
	}

	runtimes[pid] = r

	return true
}
