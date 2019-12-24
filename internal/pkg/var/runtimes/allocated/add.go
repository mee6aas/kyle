package runtimes

import (
	"github.com/mee6aas/kyle/internal/pkg/runtime"
	"github.com/pkg/errors"
)

// Add adds a runtime into the collection with specified activity name as the key.
func Add(actName string, r *runtime.Runtime) (e error) {
	if !r.IsConnected() {
		e = errors.New("Not connected runtime")
		return
	}

	rs, ok := runtimes[actName]
	if !ok {
		rs = make([]*runtime.Runtime, 0, 1)
	}

	rs = append(rs, r)
	runtimes[actName] = rs

	return
}
