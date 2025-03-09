package cli

import (
	"fmt"
	"os"

	"github.com/SailfinIO/agent/pkg/agent"
	"github.com/SailfinIO/agent/pkg/config"
	"github.com/SailfinIO/agent/pkg/utils"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for the agent CLI.
func NewRootCmd() *cobra.Command {
	// Create a logger instance for the CLI.
	logger := utils.New().WithContext("cli")

	// Load configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config: " + err.Error())
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "agent",
		Short: "Sailfin agent collects server metrics",
	}

	// 'run' command to start the agent service (daemon mode).
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the agent service",
		Run: func(cmd *cobra.Command, args []string) {
			a, err := agent.NewAgent(cfg)
			if err != nil {
				logger.Error("Error initializing agent: " + err.Error())
				os.Exit(1)
			}
			logger.Info("Starting agent service")
			if err := a.Start(); err != nil {
				logger.Error("Error running agent: " + err.Error())
				os.Exit(1)
			}
		},
	}

	// 'metrics' command to collect metrics on demand.
	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Collect and display metrics",
		Run: func(cmd *cobra.Command, args []string) {
			a, err := agent.NewAgent(cfg)
			if err != nil {
				logger.Error("Error initializing agent: " + err.Error())
				os.Exit(1)
			}
			m, err := a.CollectMetrics()
			if err != nil {
				logger.Error("Error collecting metrics: " + err.Error())
				os.Exit(1)
			}
			logger.Info(fmt.Sprintf("Metrics: %+v", m))
		},
	}

	rootCmd.AddCommand(runCmd, metricsCmd)
	return rootCmd
}
