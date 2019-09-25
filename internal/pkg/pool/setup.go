package pool

// Config holds the configuration for the pool.
type Config struct {
	AgentHost string
	AgentPort uint16
}

// Setup sets the configuration for the pool.
func Setup(conf Config) (e error) {
	rConf.AgentHost = conf.AgentHost
	rConf.AgentPort = conf.AgentPort

	return
}
