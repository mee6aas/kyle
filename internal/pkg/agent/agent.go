package agent

import (
	"os"

	"github.com/mee6aas/zeep/api"
)

var (
	localAgntHost string
	localAgntPort string
	localAgntAddr string

	serveHost string
	servePort string
	serveAddr string
)

func init() {
	localAgntHost = os.Getenv(api.AgentHostEnvKey)
	localAgntPort = os.Getenv(api.AgentPortEnvKey)
	if localAgntHost != "" && localAgntPort != "" {
		localAgntAddr = localAgntHost + ":" + localAgntPort
	}
}
