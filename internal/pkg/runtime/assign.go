package runtime

import (
	"context"

	"github.com/pkg/errors"
)

// Assign passes task to operator
func (r *Runtime) Assign(ctx context.Context, task interface{}) (e error) {
	if !r.IsAllocated() {
		e = errors.New("Task operator not allocated")
		return
	}

	r.isAssigned = true

	e = r.taskAssigner.Assign(ctx, task)
	if e != nil {
		r.isAssigned = false
		return
	}

	return
}
