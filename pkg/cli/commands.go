package cli

import (
	"fmt"
	"log"

	"github.com/SailfinIO/agent/pkg/agent"
	"github.com/SailfinIO/agent/pkg/config"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for the agent CLI.
func NewRootCmd() *cobra.Command {
	// Load configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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
				log.Fatalf("Error initializing agent: %v", err)
			}
			if err := a.Start(); err != nil {
				log.Fatalf("Error running agent: %v", err)
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
				log.Fatalf("Error initializing agent: %v", err)
			}
			m, err := a.CollectMetrics()
			if err != nil {
				log.Fatalf("Error collecting metrics: %v", err)
			}
			fmt.Printf("Metrics: %+v\n", m)
		},
	}

	rootCmd.AddCommand(runCmd, metricsCmd)
	return rootCmd
}
