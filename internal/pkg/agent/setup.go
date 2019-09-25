package agent

import (
	"strconv"

	"github.com/mee6aas/kyle/internal/pkg/pool"
)

// Config holds the configuration for the agent.
type Config struct {
	Host string
	Port uint16
}

// Setup initializes the agent.
func Setup(conf Config) (e error) {
	serveHost = conf.Host
	servePort = strconv.Itoa(int(conf.Port))
	serveAddr = serveHost + ":" + servePort

	pool.Setup(pool.Config{
		AgentHost: conf.Host,
		AgentPort: conf.Port,
	})

	return
}
