package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/SailfinIO/agent/gen/agentconfig"
)

type Config = agentconfig.AgentConfig
type RemoteHost = agentconfig.RemoteHost

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sailfin", "AgentConfig.pkl"), nil
}

// SaveConfig saves configuration to a Pkl file in the user's ~/.sailfin directory.
func SaveConfig(cfg *Config) error {
	target, err := getConfigPath()
	if err != nil {
		return err
	}
	data, err := agentconfig.Marshal(cfg)
	if err != nil {
		return err
	}
	// Ensure the directory exists.
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	return os.WriteFile(target, data, 0644)
}

// LoadConfig loads configuration from a Pkl file in the user's ~/.sailfin directory.
func LoadConfig() (*Config, error) {
	target, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	return agentconfig.LoadFromPath(context.Background(), target)
}
