package runtimes

import "github.com/mee6aas/kyle/internal/pkg/runtime"

// Take withdraws a runtime from the collection
func Take() (r *runtime.Runtime, ok bool) {
	if len(runtimes) == 0 {
		return nil, false
	}

	for k, v := range runtimes {
		delete(runtimes, k)
		return v, true
	}

	return nil, false
}
