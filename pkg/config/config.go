// pkg/confing/config.go

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the configuration for the agent.
type Config struct {
	ServerAddress string `mapstructure:"server_address"`
	APIKey        string `mapstructure:"api_key"`
	// Add additional fields as necessary.
}

// LoadConfig reads configuration from file/environment variables.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("sailfin")  // Name of config file (without extension)
	viper.SetConfigType("yaml")     // Config file format
	viper.AddConfigPath(".")        // Look in the current directory
	viper.AddConfigPath("./config") // Look for config in the config directory
	viper.AddConfigPath("$HOME")    // Look in the home directory

	// Set defaults
	viper.SetDefault("server_address", "localhost:8080")
	viper.SetDefault("api_key", "default-key")

	// Read configuration.
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	return &cfg, nil
}
