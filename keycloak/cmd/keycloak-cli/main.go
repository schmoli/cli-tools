// keycloak/cmd/keycloak-cli/main.go
package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var version = "dev"

var (
	flagURL          string
	flagRealm        string
	flagClientID     string
	flagClientSecret string
	flagInsecure     bool
	flagTargetRealm  string
)

var rootCmd = &cobra.Command{
	Use:     "keycloak-cli",
	Short:   "CLI for Keycloak Admin API",
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Keycloak URL (or KEYCLOAK_URL)")
	rootCmd.PersistentFlags().StringVar(&flagRealm, "realm", "", "Auth realm (or KEYCLOAK_REALM)")
	rootCmd.PersistentFlags().StringVar(&flagClientID, "client-id", "", "Client ID (or KEYCLOAK_CLIENT_ID)")
	rootCmd.PersistentFlags().StringVar(&flagClientSecret, "client-secret", "", "Client secret (or KEYCLOAK_CLIENT_SECRET)")
	rootCmd.PersistentFlags().BoolVarP(&flagInsecure, "insecure", "k", false, "Skip TLS verification")
	rootCmd.PersistentFlags().StringVar(&flagTargetRealm, "target-realm", "", "Target realm for queries (or KEYCLOAK_TARGET_REALM)")

	rootCmd.AddCommand(realmsCmd)
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(clientsCmd)
	rootCmd.AddCommand(rolesCmd)
	rootCmd.AddCommand(groupsCmd)
}

func envOrFlag(flag, env string) string {
	if flag != "" {
		return flag
	}
	return os.Getenv(env)
}

func getClient() (*keycloak.Client, error) {
	cfg := keycloak.Config{
		URL:          envOrFlag(flagURL, "KEYCLOAK_URL"),
		Realm:        envOrFlag(flagRealm, "KEYCLOAK_REALM"),
		ClientID:     envOrFlag(flagClientID, "KEYCLOAK_CLIENT_ID"),
		ClientSecret: envOrFlag(flagClientSecret, "KEYCLOAK_CLIENT_SECRET"),
		Insecure:     flagInsecure,
	}
	return keycloak.NewClient(cfg)
}

func getTargetRealm() string {
	return envOrFlag(flagTargetRealm, "KEYCLOAK_TARGET_REALM")
}

func handleError(err error) {
	keycloak.PrintError(err)
	if ke, ok := err.(*keycloak.KeycloakError); ok {
		os.Exit(ke.ExitCode())
	}
	os.Exit(1)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
