package pool

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mee6aas/kyle/internal/pkg/runtime"
	runtimesConnected "github.com/mee6aas/kyle/internal/pkg/var/runtimes/connected"
	runtimesPended "github.com/mee6aas/kyle/internal/pkg/var/runtimes/pended"
)

func spawn(ctx context.Context, conf runtime.Config) (e error) {
	r, e := runtime.NewRuntime(conf)
	if e != nil {
		e = errors.Wrap(e, "Failed to create a new runtime")
		return
	}

	defer func() {
		if e == nil {
			return
		}

		r.Cancel()
	}()

	e = r.Start()
	if e != nil {
		e = errors.Wrap(e, "Failed to start the runtime")
		return
	}

	pid, ok := r.PID()
	if !ok {
		e = errors.New("Failed to get PID from the runtime")
		return
	}

	onReleased, ok := runtimesPended.Add(r)
	if !ok {
		e = errors.New("Failed to add the runtime to the collection for the pended runtimes")
		return
	}
	log.WithField("pid", pid).Debug("Runtime pended")

	select {
	case <-ctx.Done():
		e = ctx.Err()
		return
	case <-onReleased:
		if ok := runtimesConnected.Add(r); !ok {
			e = errors.New("Failed to add the runtime to the collection for the allocated runtimes")
			return
		}
		log.WithField("pid", pid).Debug("Runtime connected")
		return
	}
}
