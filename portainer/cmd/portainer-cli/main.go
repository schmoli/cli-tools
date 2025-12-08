package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/portainer/pkg/portainer"
)

var version = "dev"

var (
	flagURL      string
	flagToken    string
	flagInsecure bool
)

var rootCmd = &cobra.Command{
	Use:     "portainer-cli",
	Short:   "CLI for Portainer API",
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Portainer URL (or set PORTAINER_URL)")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "API token (or set PORTAINER_TOKEN)")
	rootCmd.PersistentFlags().BoolVarP(&flagInsecure, "insecure", "k", false, "Skip TLS certificate verification")

	rootCmd.AddCommand(stacksCmd)
	rootCmd.AddCommand(endpointsCmd)
	rootCmd.AddCommand(containersCmd)
}

func parseID(arg string) (int64, error) {
	var id int64
	if _, err := fmt.Sscanf(arg, "%d", &id); err != nil {
		return 0, portainer.ConfigError(fmt.Sprintf("invalid ID: %s", arg))
	}
	if id <= 0 {
		return 0, portainer.ConfigError("ID must be positive")
	}
	return id, nil
}

func getConfig() (string, string, error) {
	url := flagURL
	if url == "" {
		url = os.Getenv("PORTAINER_URL")
	}
	if url == "" {
		return "", "", portainer.ConfigError("missing URL. Use --url or set PORTAINER_URL")
	}

	token := flagToken
	if token == "" {
		token = os.Getenv("PORTAINER_TOKEN")
	}
	if token == "" {
		return "", "", portainer.ConfigError("missing token. Use --token or set PORTAINER_TOKEN")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", "", portainer.ConfigError("URL must start with http:// or https://")
	}

	return url, token, nil
}

func getClient() (*portainer.Client, error) {
	url, token, err := getConfig()
	if err != nil {
		return nil, err
	}
	return portainer.NewClient(url, token, flagInsecure), nil
}

func handleError(err error) {
	portainer.PrintError(err)
	if pe, ok := err.(*portainer.PortainerError); ok {
		os.Exit(pe.ExitCode())
	}
	os.Exit(1)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
