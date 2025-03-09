package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/SailfinIO/agent/pkg/agent"
	"github.com/SailfinIO/agent/pkg/config"
	"github.com/spf13/cobra"
)

func main() {
	// Load configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "agent",
		Short: "Sailfin agent collects server metrics",
	}

	// Command to run the agent as a service.
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the agent service",
		Run: func(cmd *cobra.Command, args []string) {
			a, err := agent.NewAgent(cfg)
			if err != nil {
				log.Fatalf("Error initializing agent: %v", err)
			}
			if err := a.Start(); err != nil {
				log.Fatalf("Error starting agent: %v", err)
			}
		},
	}

	// Command to retrieve stored metrics via CLI.
	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Retrieve stored metrics snapshots",
		Run: func(cmd *cobra.Command, args []string) {
			a, err := agent.NewAgent(cfg)
			if err != nil {
				log.Fatalf("Error initializing agent: %v", err)
			}

			limit, _ := cmd.Flags().GetInt("limit")
			fromStr, _ := cmd.Flags().GetString("from")
			toStr, _ := cmd.Flags().GetString("to")

			// If both "from" and "to" flags are provided, query by time range.
			if fromStr != "" && toStr != "" {
				fromUnix, err1 := strconv.ParseInt(fromStr, 10, 64)
				toUnix, err2 := strconv.ParseInt(toStr, 10, 64)
				if err1 != nil || err2 != nil {
					log.Fatalf("Invalid from/to timestamps provided")
				}
				from := time.Unix(fromUnix, 0)
				to := time.Unix(toUnix, 0)
				snaps, err := a.GetSnapshotsByTime(from, to)
				if err != nil {
					log.Fatalf("Error retrieving snapshots: %v", err)
				}
				log.Printf("Snapshots (from %v to %v): %+v", from, to, snaps)
				return
			}

			// If "limit" is provided, return the latest N snapshots.
			if limit > 0 {
				snaps, err := a.GetSnapshotsByLimit(limit)
				if err != nil {
					log.Fatalf("Error retrieving snapshots: %v", err)
				}
				log.Printf("Latest %d snapshots: %+v", limit, snaps)
				return
			}

			// Default: return the latest snapshot.
			snap, err := a.GetLatestSnapshot()
			if err != nil {
				log.Fatalf("Error retrieving the latest snapshot: %v", err)
			}
			log.Printf("Latest Snapshot: %+v", snap)
		},
	}

	// Define flags for the metrics command.
	metricsCmd.Flags().Int("limit", 0, "Number of latest snapshots to retrieve")
	metricsCmd.Flags().String("from", "", "Unix timestamp start for snapshot query")
	metricsCmd.Flags().String("to", "", "Unix timestamp end for snapshot query")

	rootCmd.AddCommand(runCmd, metricsCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
