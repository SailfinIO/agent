// pkg/cli/commands.go
package cli

import (
	"os"

	"github.com/SailfinIO/agent/pkg/config"
	"github.com/SailfinIO/agent/pkg/utils"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for the CLI.
func NewRootCmd() *cobra.Command {
	logger := utils.New().WithContext("cli")

	// Load configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config: " + err.Error())
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "sailfin",
		Short: "Sailfin collects server metrics",
	}

	// Add agent commands.
	rootCmd.AddCommand(newAgentCmd(cfg))
	// Add remote commands.
	rootCmd.AddCommand(newRemoteCmd())

	return rootCmd
}
