package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/toli/portainer-cli/pkg/portainer"
)

var (
	flagURL   string
	flagToken string
)

var rootCmd = &cobra.Command{
	Use:     "portainer-cli",
	Short:   "CLI for Portainer API",
	Version: "0.1.0",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Portainer URL (or set PORTAINER_URL)")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "API token (or set PORTAINER_TOKEN)")

	rootCmd.AddCommand(stacksCmd)
	rootCmd.AddCommand(endpointsCmd)
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

	return url, token, nil
}

func getClient() (*portainer.Client, error) {
	url, token, err := getConfig()
	if err != nil {
		return nil, err
	}
	return portainer.NewClient(url, token), nil
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
