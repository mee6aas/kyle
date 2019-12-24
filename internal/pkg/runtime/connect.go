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

// Connect sets the specified task assigner to this runtime.
func (r *Runtime) Connect(ta TaskAssigner) (e error) {
	if r.isConnected {
		e = errors.New("Already allocated worker")
		return
	}

	r.isConnected = true
	r.taskAssigner = ta

	return
}
