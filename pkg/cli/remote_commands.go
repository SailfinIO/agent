package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/SailfinIO/agent/gen/agentconfig"
	"github.com/SailfinIO/agent/pkg/config"
	"github.com/spf13/cobra"
)

// saveConfig saves the updated AgentConfig to disk.
func saveConfig(cfg *config.Config) error {
	target := filepath.Join("pkl", "AgentConfig.pkl")
	data, err := agentconfig.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(target, data, 0644)
}

// newRemoteCmd creates the "remote" command group.
func newRemoteCmd() *cobra.Command {
	remoteCmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage remote server configurations",
	}

	// "add" command to add a new remote host.
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a remote server",
		Run: func(cmd *cobra.Command, args []string) {
			host, _ := cmd.Flags().GetString("host")
			userFlag, _ := cmd.Flags().GetString("user")
			password, _ := cmd.Flags().GetString("password")
			privateKey, _ := cmd.Flags().GetString("private-key")

			// If no host is provided, require that the user be explicit.
			if host == "" || userFlag == "" {
				fmt.Println("Host and user are required.")
				os.Exit(1)
			}
			if password != "" && privateKey != "" {
				fmt.Println("Please provide only one authentication method: either password or private key.")
				os.Exit(1)
			}
			if password == "" && privateKey == "" {
				fmt.Println("No authentication method provided. Please provide either a password or a private key.")
				os.Exit(1)
			}

			// Load current configuration.
			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				os.Exit(1)
			}

			// Check if remote host already exists.
			for _, r := range cfg.RemoteHosts {
				if r.Host == host && r.User == userFlag {
					fmt.Println("Remote host already exists in configuration.")
					os.Exit(1)
				}
			}

			// Create a new RemoteHost.
			newRemote := &agentconfig.RemoteHost{
				Host:       host,
				User:       userFlag,
				Password:   nil,
				PrivateKey: nil,
			}
			if password != "" {
				newRemote.Password = &password
			}
			if privateKey != "" {
				newRemote.PrivateKey = &privateKey
			}

			// Append the new remote host.
			cfg.RemoteHosts = append(cfg.RemoteHosts, newRemote)

			// Save the updated configuration.
			if err := saveConfig(cfg); err != nil {
				fmt.Printf("Error saving updated configuration: %v\n", err)
				os.Exit(1)
			}

			authType := "password"
			if privateKey != "" {
				authType = "private key"
			}
			fmt.Printf("Remote host %s for user %s added successfully using %s authentication.\n", host, userFlag, authType)
		},
	}
	addCmd.Flags().String("host", "", "Hostname or IP address of the remote server")
	addCmd.Flags().String("user", "", "Username for remote authentication")
	addCmd.Flags().String("password", "", "Password for remote authentication")
	addCmd.Flags().String("private-key", "", "Path to the private key for remote authentication")

	// "set" command to update an existing remote host.
	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Update an existing remote host configuration",
		Run: func(cmd *cobra.Command, args []string) {
			host, _ := cmd.Flags().GetString("host")
			userFlag, _ := cmd.Flags().GetString("user")
			password, _ := cmd.Flags().GetString("password")
			privateKey, _ := cmd.Flags().GetString("private-key")

			if host == "" || userFlag == "" {
				fmt.Println("Host and user are required to identify the remote host to update.")
				os.Exit(1)
			}
			if password != "" && privateKey != "" {
				fmt.Println("Please provide only one authentication method: either password or private key.")
				os.Exit(1)
			}
			if password == "" && privateKey == "" {
				fmt.Println("No new authentication method provided. Provide either a password or a private key to update.")
				os.Exit(1)
			}

			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				os.Exit(1)
			}

			updated := false
			for _, r := range cfg.RemoteHosts {
				if r.Host == host && r.User == userFlag {
					if password != "" {
						r.Password = &password
						r.PrivateKey = nil
					} else {
						r.PrivateKey = &privateKey
						r.Password = nil
					}
					updated = true
					break
				}
			}

			if !updated {
				fmt.Println("Remote host not found in configuration.")
				os.Exit(1)
			}

			if err := saveConfig(cfg); err != nil {
				fmt.Printf("Error saving updated configuration: %v\n", err)
				os.Exit(1)
			}

			authType := "password"
			if privateKey != "" {
				authType = "private key"
			}
			fmt.Printf("Remote host %s for user %s updated successfully to use %s authentication.\n", host, userFlag, authType)
		},
	}
	setCmd.Flags().String("host", "", "Hostname or IP address of the remote server to update")
	setCmd.Flags().String("user", "", "Username for remote authentication")
	setCmd.Flags().String("password", "", "New password for remote authentication")
	setCmd.Flags().String("private-key", "", "New path to the private key for remote authentication")

	// "list" command to display all remote hosts.
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all configured remote hosts",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				os.Exit(1)
			}
			if len(cfg.RemoteHosts) == 0 {
				fmt.Println("No remote hosts configured.")
				return
			}
			fmt.Println("Configured remote hosts:")
			for _, r := range cfg.RemoteHosts {
				authType := "password"
				if r.PrivateKey != nil {
					authType = "private key"
				}
				fmt.Printf("Host: %s, User: %s, Auth: %s\n", r.Host, r.User, authType)
			}
		},
	}

	// "delete" command to remove a remote host.
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a remote host configuration",
		Run: func(cmd *cobra.Command, args []string) {
			host, _ := cmd.Flags().GetString("host")
			userFlag, _ := cmd.Flags().GetString("user")

			if host == "" || userFlag == "" {
				fmt.Println("Host and user are required to identify the remote host to delete.")
				os.Exit(1)
			}

			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				os.Exit(1)
			}

			// Filter out the remote host to be deleted.
			newRemotes := []*agentconfig.RemoteHost{}
			deleted := false
			for _, r := range cfg.RemoteHosts {
				if r.Host == host && r.User == userFlag {
					deleted = true
					continue
				}
				newRemotes = append(newRemotes, r)
			}

			if !deleted {
				fmt.Println("Remote host not found in configuration.")
				os.Exit(1)
			}

			cfg.RemoteHosts = newRemotes

			if err := saveConfig(cfg); err != nil {
				fmt.Printf("Error saving updated configuration: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Remote host %s for user %s deleted successfully.\n", host, userFlag)
		},
	}
	deleteCmd.Flags().String("host", "", "Hostname or IP address of the remote server to delete")
	deleteCmd.Flags().String("user", "", "Username for remote authentication of the host to delete")

	// "local" command to add the local machine as a remote host by default.
	localCmd := &cobra.Command{
		Use:   "local",
		Short: "Add local machine as a remote host (defaults to 'localhost' and current OS user)",
		Run: func(cmd *cobra.Command, args []string) {
			// Get current OS user.
			currentUser, err := user.Current()
			if err != nil {
				fmt.Printf("Error fetching current user: %v\n", err)
				os.Exit(1)
			}
			host := "localhost"
			username := currentUser.Username

			// Load current configuration.
			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				os.Exit(1)
			}

			// Check if local host already exists.
			for _, r := range cfg.RemoteHosts {
				if r.Host == host && r.User == username {
					fmt.Println("Local host is already configured.")
					return
				}
			}

			// Create a new RemoteHost entry for the local machine.
			localRemote := &agentconfig.RemoteHost{
				Host:       host,
				User:       username,
				Password:   nil,
				PrivateKey: nil,
			}
			// (Optionally, you could decide on an authentication method if needed.)

			cfg.RemoteHosts = append(cfg.RemoteHosts, localRemote)

			if err := saveConfig(cfg); err != nil {
				fmt.Printf("Error saving updated configuration: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Local host (%s@%s) added successfully.\n", username, host)
		},
	}

	// "install" command to install the agent on a remote server.
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install the agent on a remote server",
		Run: func(cmd *cobra.Command, args []string) {
			host, _ := cmd.Flags().GetString("host")
			userFlag, _ := cmd.Flags().GetString("user")
			privateKey, _ := cmd.Flags().GetString("private-key")

			if host == "" || userFlag == "" {
				fmt.Println("Remote host and user are required.")
				os.Exit(1)
			}

			sshArgs := []string{}
			if privateKey != "" {
				sshArgs = append(sshArgs, "-i", privateKey)
			}
			target := fmt.Sprintf("%s@%s", userFlag, host)
			sshArgs = append(sshArgs, target)
			remoteCmdStr := "bash -c 'curl -sL https://raw.githubusercontent.com/SailfinIO/agent/main/scripts/install.sh | bash'"
			sshArgs = append(sshArgs, remoteCmdStr)

			fmt.Printf("Installing agent on remote host %s as user %s...\n", host, userFlag)
			out, err := exec.Command("ssh", sshArgs...).CombinedOutput()
			if err != nil {
				fmt.Printf("Error installing agent: %v\nOutput: %s\n", err, string(out))
				os.Exit(1)
			}
			fmt.Printf("Installation successful on remote host. Output:\n%s\n", string(out))
		},
	}
	installCmd.Flags().String("host", "", "Hostname or IP address of the remote server")
	installCmd.Flags().String("user", "", "Username for remote SSH authentication")
	installCmd.Flags().String("private-key", "", "Path to the SSH private key for authentication")

	remoteCmd.AddCommand(addCmd, setCmd, listCmd, deleteCmd, localCmd, installCmd)
	return remoteCmd
}
