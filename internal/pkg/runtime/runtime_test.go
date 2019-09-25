package runtime

import (
	"bufio"
	"strconv"
	"testing"

	"gotest.tools/assert"

	"github.com/mee6aas/zeep/api"
)

func TestNewRuntime(t *testing.T) {
	const (
		testHost        = "Zeep"
		testPort uint16 = 5122
	)

	runtimeExePath = "./testdata/printEnv.sh"

	r, e := NewRuntime(Config{
		AgentHost: testHost,
		AgentPort: testPort,
	})
	assert.NilError(t, e, "Failed to create a new Runtime")
	defer r.Cancel()

	r.cmd.Args = []string{
		"",
		api.AgentHostEnvKey,
		api.AgentPortEnvKey,
	}

	stdout, e := r.cmd.StdoutPipe()
	assert.NilError(t, e, "Failed to get stdout from the runtime")

	reader := bufio.NewReader(stdout)

	e = r.Start()
	assert.NilError(t, e, "Failed to start the runtime")

	l, e := reader.ReadString('\n')
	assert.NilError(t, e, "Failed to read line from the runtime")
	assert.Equal(t, l, testHost+"\n")

	l, e = reader.ReadString('\n')
	assert.NilError(t, e, "Failed to read line from the runtime")
	assert.Equal(t, l, strconv.Itoa(int(testPort))+"\n")

}
