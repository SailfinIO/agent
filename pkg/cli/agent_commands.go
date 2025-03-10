// pkg/cli/agent_commands.go
package cli

import (
	"fmt"
	"os"

	"github.com/SailfinIO/agent/pkg/agent"
	"github.com/SailfinIO/agent/pkg/config"
	"github.com/SailfinIO/agent/pkg/utils"
	"github.com/spf13/cobra"
)

// newAgentCmd creates the "agent" command group.
func newAgentCmd(cfg *config.Config) *cobra.Command {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage the Sailfin agent",
	}

	// "start" command to run the agent as a daemon.
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the agent service",
		Run: func(cmd *cobra.Command, args []string) {
			logger := utils.New().WithContext("agent")
			a, err := agent.NewAgent(cfg)
			if err != nil {
				logger.Error("Error initializing agent: " + err.Error())
				os.Exit(1)
			}
			logger.Info("Starting agent service")
			if err := a.Start(); err != nil {
				logger.Error("Error starting agent: " + err.Error())
				os.Exit(1)
			}
		},
	}

	// "stop" command to stop the agent service.
	// Here you might want to implement graceful shutdown via signals or a PID file.
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the agent service",
		Run: func(cmd *cobra.Command, args []string) {
			// This is a placeholder. You might read a PID file and send a termination signal.
			fmt.Println("Stopping the agent service...")
			// For example:
			// pid, err := ioutil.ReadFile("agent.pid")
			// if err != nil { ... }
			// syscall.Kill(pid, syscall.SIGTERM)
		},
	}

	// "metrics" command to retrieve metrics.
	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Retrieve stored metrics snapshots",
		Run: func(cmd *cobra.Command, args []string) {
			logger := utils.New().WithContext("agent")
			a, err := agent.NewAgent(cfg)
			if err != nil {
				logger.Error("Error initializing agent: " + err.Error())
				os.Exit(1)
			}

			limit, _ := cmd.Flags().GetInt("limit")
			fromStr, _ := cmd.Flags().GetString("from")
			toStr, _ := cmd.Flags().GetString("to")

			// Time-range query if both "from" and "to" flags are provided.
			if fromStr != "" && toStr != "" {
				// (Timestamp parsing logic, as in your existing code.)
				// ...
				logger.Info("Time-range metrics not yet implemented")
				return
			}

			if limit > 0 {
				snaps, err := a.GetSnapshotsByLimit(limit)
				if err != nil {
					logger.Error("Error retrieving snapshots: " + err.Error())
					os.Exit(1)
				}
				logger.Info(fmt.Sprintf("Latest %d snapshots: %+v", limit, snaps))
				return
			}

			snap, err := a.GetLatestSnapshot()
			if err != nil {
				logger.Error("Error retrieving the latest snapshot: " + err.Error())
				os.Exit(1)
			}
			logger.Info(fmt.Sprintf("Latest Snapshot: %+v", snap))
		},
	}
	metricsCmd.Flags().Int("limit", 0, "Number of latest snapshots to retrieve")
	metricsCmd.Flags().String("from", "", "Unix timestamp start for snapshot query")
	metricsCmd.Flags().String("to", "", "Unix timestamp end for snapshot query")

	agentCmd.AddCommand(startCmd, stopCmd, metricsCmd)
	return agentCmd
}
