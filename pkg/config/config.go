// pkg/config/config.go
package config

import (
	"context"

	"github.com/SailfinIO/agent/gen/agentconfig"
)

// You can alias the generated types so that other parts of your project can use them.
type Config = agentconfig.AgentConfig
type RemoteHost = agentconfig.RemoteHost

// LoadConfig loads configuration from a Pkl file.
func LoadConfig() (*Config, error) {
	// You might want to allow the config file path to be dynamic; for now, we hard-code it.
	return agentconfig.LoadFromPath(context.Background(), "pkl/AgentConfig.pkl")
}
