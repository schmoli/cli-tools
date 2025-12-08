package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/pve/pkg/pve"
)

var version = "dev"

var (
	flagURL         string
	flagTokenID     string
	flagTokenSecret string
	flagInsecure    bool
)

var rootCmd = &cobra.Command{
	Use:     "pve-cli",
	Short:   "CLI for Proxmox VE API",
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Proxmox URL (or set PVE_URL)")
	rootCmd.PersistentFlags().StringVar(&flagTokenID, "token-id", "", "Token ID (or set PVE_TOKEN_ID)")
	rootCmd.PersistentFlags().StringVar(&flagTokenSecret, "token", "", "Token secret (or set PVE_TOKEN_SECRET)")
	rootCmd.PersistentFlags().BoolVarP(&flagInsecure, "insecure", "k", false, "Skip TLS certificate verification")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
}

func getConfig() (string, string, string, error) {
	url := flagURL
	if url == "" {
		url = os.Getenv("PVE_URL")
	}
	if url == "" {
		return "", "", "", pve.ConfigError("missing URL. Use --url or set PVE_URL")
	}

	tokenID := flagTokenID
	if tokenID == "" {
		tokenID = os.Getenv("PVE_TOKEN_ID")
	}
	if tokenID == "" {
		return "", "", "", pve.ConfigError("missing token ID. Use --token-id or set PVE_TOKEN_ID")
	}

	tokenSecret := flagTokenSecret
	if tokenSecret == "" {
		tokenSecret = os.Getenv("PVE_TOKEN_SECRET")
	}
	if tokenSecret == "" {
		return "", "", "", pve.ConfigError("missing token secret. Use --token or set PVE_TOKEN_SECRET")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", "", "", pve.ConfigError("URL must start with http:// or https://")
	}

	return url, tokenID, tokenSecret, nil
}

func getClient() (*pve.Client, error) {
	url, tokenID, tokenSecret, err := getConfig()
	if err != nil {
		return nil, err
	}
	return pve.NewClient(url, tokenID, tokenSecret, flagInsecure), nil
}

func parseVMID(arg string) (int64, error) {
	var id int64
	if _, err := fmt.Sscanf(arg, "%d", &id); err != nil {
		return 0, pve.ConfigError(fmt.Sprintf("invalid VMID: %s", arg))
	}
	if id <= 0 {
		return 0, pve.ConfigError("VMID must be positive")
	}
	return id, nil
}

func handleError(err error) {
	pve.PrintError(err)
	if pe, ok := err.(*pve.PveError); ok {
		os.Exit(pe.ExitCode())
	}
	os.Exit(1)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VMs and LXCs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return nil
		}

		guests, err := client.ListGuests()
		if err != nil {
			handleError(err)
			return nil
		}

		if err := pve.PrintYAML(guests); err != nil {
			handleError(err)
		}
		return nil
	},
}

var startCmd = &cobra.Command{
	Use:   "start <vmid>",
	Short: "Start a VM or LXC",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vmid, err := parseVMID(args[0])
		if err != nil {
			handleError(err)
			return nil
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return nil
		}

		vmType, name, err := client.FindGuestType(vmid)
		if err != nil {
			handleError(err)
			return nil
		}

		if err := client.StartGuest(vmid, vmType); err != nil {
			handleError(err)
			return nil
		}

		result := pve.ActionResult{
			VMID:   vmid,
			Name:   name,
			Action: "started",
		}
		if err := pve.PrintYAML(result); err != nil {
			handleError(err)
		}
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop <vmid>",
	Short: "Stop a VM or LXC",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vmid, err := parseVMID(args[0])
		if err != nil {
			handleError(err)
			return nil
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return nil
		}

		vmType, name, err := client.FindGuestType(vmid)
		if err != nil {
			handleError(err)
			return nil
		}

		if err := client.StopGuest(vmid, vmType); err != nil {
			handleError(err)
			return nil
		}

		result := pve.ActionResult{
			VMID:   vmid,
			Name:   name,
			Action: "stopped",
		}
		if err := pve.PrintYAML(result); err != nil {
			handleError(err)
		}
		return nil
	},
}
