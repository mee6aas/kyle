package runtime

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strconv"

	"github.com/mee6aas/zeep/api"
)

var (
	runtimeExePath string
)

// Runtime holds process for the runtime.
type Runtime struct {
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc

	isAllocated  bool
	taskAssigner TaskAssigner
	isAssigned   bool
}

// Cancel expires the context included in the runtime.
func (r *Runtime) Cancel() { r.cancel() }

// Start start the runtime.
func (r *Runtime) Start() (e error) {
	stdout, e := r.cmd.StdoutPipe()
	if e != nil {
		return
	}
	stderr, e := r.cmd.StderrPipe()
	if e != nil {
		return
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stdout, stderr)

	if e = r.cmd.Start(); e != nil {
		return
	}

	return
}

// PID returns PID of the runtime.
// Notice that the NewRuntime executes the spawner of the runtime process.
// So that the PID of the actual process may be different from this PID.
func (r *Runtime) PID() (int, bool) {
	if r.cmd.Process == nil {
		return 0, false
	}

	return r.cmd.Process.Pid, true
}

// IsAllocated checks if the task is allocated to this runtime.
func (r *Runtime) IsAllocated() bool { return r.isAllocated }

// IsAssigned checks if the task is assigned to this runtime.
func (r *Runtime) IsAssigned() bool { return r.isAssigned }

// Resolve set isAssigned flag to false
func (r *Runtime) Resolve() { r.isAssigned = false }

// Config holds the configuration for the runtime process.
type Config struct {
	AgentHost string
	AgentPort uint16
}

// NewRuntime creates a new runtime process based on the given configuration.
func NewRuntime(conf Config) (r *Runtime, e error) {
	r = &Runtime{}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.cmd = exec.CommandContext(r.ctx, runtimeExePath)

	// set env
	{
		env := os.Environ()
		env = append(env,
			api.AgentHostEnvKey+"="+conf.AgentHost,
			api.AgentPortEnvKey+"="+strconv.Itoa(int(conf.AgentPort)),
		)
		r.cmd.Env = env
	}

	return
}

func init() {
	runtimeExePath = api.RuntimeSpawn
}
