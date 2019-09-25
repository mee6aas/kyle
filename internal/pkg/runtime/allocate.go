package runtime

import (
	"context"

	"github.com/pkg/errors"
)

// TaskAssigner passes task to its task operator.
type TaskAssigner interface {
	Assign(context.Context, interface{}) error
	Close()
}

// Allocate sets the specified task assigner to this runtime.
func (r *Runtime) Allocate(ta TaskAssigner) (e error) {
	if r.isAllocated {
		e = errors.New("Already allocated worker")
		return
	}

	r.isAllocated = true
	r.taskAssigner = ta

	return
}
