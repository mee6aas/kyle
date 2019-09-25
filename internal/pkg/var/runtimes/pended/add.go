package runtimes

import (
	"github.com/mee6aas/kyle/internal/pkg/runtime"
)

// Add adds a runtime into the collection.
// The PID of the runtime used as a key.
// The returned channel notifies the specified runtime is released.
// It returns `false` if the runtime is already added or fails to get PID of runtime.
func Add(r *runtime.Runtime) (<-chan struct{}, bool) {
	pid, ok := r.PID()
	if !ok {
		// runtime not started
		return nil, false
	}

	if _, ok = runtimes[pid]; ok {
		return nil, false
	}

	var onRelease = make(chan struct{}, 1)

	runtimes[pid] = elem{
		runtime:   r,
		onRelease: onRelease,
	}

	return onRelease, true
}
