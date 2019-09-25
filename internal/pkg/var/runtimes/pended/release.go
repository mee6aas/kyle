package runtimes

import (
	"github.com/mee6aas/kyle/internal/pkg/runtime"
)

// Release releases the hanged runtime with the specified PID in the collection
func Release(pid int, ta runtime.TaskAssigner) bool {
	r, ok := runtimes[pid]
	if !ok {
		return false
	}
	delete(runtimes, pid)

	if e := r.runtime.Allocate(ta); e != nil {
		// already allocated
		return false
	}

	r.onRelease <- struct{}{}

	return true
}
