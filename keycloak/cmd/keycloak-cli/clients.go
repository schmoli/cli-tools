// keycloak/cmd/keycloak-cli/clients.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "Manage clients",
}

var clientsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List clients in realm",
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		clients, err := client.ListClients(realm)
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(clients); err != nil {
			handleError(err)
		}
	},
}

var clientsGetCmd = &cobra.Command{
	Use:   "get <client-id>",
	Short: "Get client details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		cl, err := client.GetClient(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(cl); err != nil {
			handleError(err)
		}
	},
}

var clientsSessionsCmd = &cobra.Command{
	Use:   "sessions <client-uuid>",
	Short: "List client sessions",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		sessions, err := client.GetClientSessions(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(sessions); err != nil {
			handleError(err)
		}
	},
}

func init() {
	clientsCmd.AddCommand(clientsListCmd)
	clientsCmd.AddCommand(clientsGetCmd)
	clientsCmd.AddCommand(clientsSessionsCmd)
}
