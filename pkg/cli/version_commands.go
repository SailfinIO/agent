package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/SailfinIO/agent/pkg/version"
	"github.com/spf13/cobra"
)

// release represents the minimal JSON structure for a GitHub release.
type release struct {
	TagName string `json:"tag_name"`
}

// NewVersionCmd creates the "version" command group.
func NewVersionCmd() *cobra.Command {
	verCmd := &cobra.Command{
		Use:   "version",
		Short: "Manage and display CLI version information",
	}

	// Display current version.
	currentCmd := &cobra.Command{
		Use:   "current",
		Short: "Display the current version of the Sailfin CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Sailfin CLI version: %s\n", version.Version)
		},
	}

	// List installed versions.
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all installed versions",
		Run: func(cmd *cobra.Command, args []string) {
			installDir := os.Getenv("SAILFIN_INSTALL_DIR")
			if installDir == "" {
				installDir = "/usr/local/sailfin/versions"
			}
			entries, err := os.ReadDir(installDir)
			if err != nil {
				fmt.Printf("Error reading install directory: %v\n", err)
				return
			}
			fmt.Println("Installed versions:")
			for _, entry := range entries {
				if entry.IsDir() {
					fmt.Println(" -", entry.Name())
				}
			}
		},
	}

	// Set the global version (update the symlink).
	setCmd := &cobra.Command{
		Use:   "set [version]",
		Short: "Set the global active version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			versionName := args[0]
			installDir := os.Getenv("SAILFIN_INSTALL_DIR")
			if installDir == "" {
				installDir = "/usr/local/sailfin/versions"
			}
			targetPath := filepath.Join(installDir, versionName, "sailfin")
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				fmt.Printf("Version %s is not installed at %s\n", versionName, targetPath)
				return
			}
			globalLink := "/usr/local/bin/sailfin"
			if err := os.Remove(globalLink); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Error removing old symlink: %v\n", err)
				return
			}
			if err := os.Symlink(targetPath, globalLink); err != nil {
				fmt.Printf("Error creating symlink: %v\n", err)
				return
			}
			fmt.Printf("Global version set to %s\n", versionName)
		},
	}

	// Update command: downloads and installs a new version.
	updateCmd := &cobra.Command{
		Use:   "update [version]",
		Short: "Download and install a new version of the Sailfin CLI",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var newVersion string
			if len(args) == 1 {
				newVersion = args[0]
			} else {
				// Default to the version embedded in the binary (updated dynamically in your workflow)
				newVersion = version.Version
			}

			// Path to the install script; adjust if needed.
			installScript := "./scripts/install.sh"

			// Prepare the command: set the VERSION environment variable for the install script.
			c := exec.Command("bash", installScript)
			c.Env = append(os.Environ(), "VERSION="+newVersion)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			fmt.Printf("Updating Sailfin CLI to version %s...\n", newVersion)
			if err := c.Run(); err != nil {
				fmt.Printf("Update failed: %v\n", err)
				return
			}
			fmt.Printf("Successfully updated to version %s\n", newVersion)
		},
	}

	// List remote versions command.
	listRemoteCmd := &cobra.Command{
		Use:   "list-remote",
		Short: "List available remote versions from GitHub",
		Run: func(cmd *cobra.Command, args []string) {
			// GitHub API endpoint for releases.
			url := "https://api.github.com/repos/SailfinIO/agent/releases"
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Get(url)
			if err != nil {
				fmt.Printf("Error fetching releases: %v\n", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
				return
			}

			var releases []release
			if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
				fmt.Printf("Error decoding releases: %v\n", err)
				return
			}

			fmt.Println("Available remote versions:")
			for _, r := range releases {
				fmt.Println(" -", r.TagName)
			}
		},
	}

	// Add subcommands to the version command.
	verCmd.AddCommand(currentCmd, listCmd, setCmd, updateCmd, listRemoteCmd)
	return verCmd
}
