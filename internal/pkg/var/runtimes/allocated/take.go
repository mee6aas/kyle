package runtimes

import (
	"github.com/mee6aas/kyle/internal/pkg/runtime"
	"github.com/pkg/errors"
)

// Take withdraws a runtime associated with the specified activity name from the collection.
func Take(actName string) (r *runtime.Runtime, e error) {
	for {
		rs, ok := runtimes[actName]
		if !ok {
			e = errors.New("Activity name not exists")
			return
		}

		if len(rs) == 0 {
			e = errors.New("Empty collection")
			return
		}

		r, runtimes[actName] = rs[0], rs[1:]

		if r.IsConnected() {
			return
		}
	}
}
