package v1

import (
	"context"
	"net"
	"os"

	"github.com/lesomnus/go-netstat/netstat"
	ps "github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	runtimes "github.com/mee6aas/kyle/internal/pkg/var/runtimes/pended"

	v1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

// TaskAssigner holds stream for assigning a task to the worker.
type TaskAssigner struct {
	ctx    context.Context
	stream chan<- v1.Task
}

// Assign sends the specified task to the worker.
func (ta TaskAssigner) Assign(ctx context.Context, t interface{}) (e error) {
	select {
	case <-ta.ctx.Done():
		e = errors.New("Disconnected")
	case <-ctx.Done():
		e = ctx.Err()
	case ta.stream <- (t.(v1.Task)):
	}

	return
}

// Close closes the connected channel
func (ta TaskAssigner) Close() {
	close(ta.stream)
}

// Connected is invoked when invokee client connected.
func (h Handle) Connected(
	ctx context.Context,
	addr *net.TCPAddr,
	stream chan<- v1.Task,
) (e error) {
	socks, e := netstat.TCPSocks(func(entry *netstat.SockTabEntry) bool {
		return int(entry.LocalAddr.Port) == addr.Port
	})
	if e != nil {
		e = errors.Wrap(e, "Failed to find TCP socket")
		return
	}
	if len(socks) == 0 {
		e = errors.New("Not found")
		return
	}

	sock := socks[0]
	pid := sock.Process.PID

	log.WithField("pid", pid).Debug("PID decided from the address")

	for {
		p, e := ps.FindProcess(pid)
		if e != nil {
			e = errors.Wrapf(e, "Failed to find process %d", pid)
			return e
		}

		ppid := p.PPid()

		if ppid == os.Getpid() {
			break
		}

		pid = ppid
	}

	log.WithField("pid", pid).Debug("PID verified")

	if ok := runtimes.Release(pid, &TaskAssigner{
		ctx:    ctx,
		stream: stream,
	}); !ok {
		// not found or already released
		e = errors.New("Invalid connection")
		return
	}

	return
}
