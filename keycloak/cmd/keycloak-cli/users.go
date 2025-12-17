// keycloak/cmd/keycloak-cli/users.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users in realm",
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

		users, err := client.ListUsers(realm)
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(users); err != nil {
			handleError(err)
		}
	},
}

var usersGetCmd = &cobra.Command{
	Use:   "get <user-id>",
	Short: "Get user details",
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

		user, err := client.GetUser(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(user); err != nil {
			handleError(err)
		}
	},
}

var usersSessionsCmd = &cobra.Command{
	Use:   "sessions <user-id>",
	Short: "List user sessions",
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

		sessions, err := client.GetUserSessions(realm, args[0])
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
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersGetCmd)
	usersCmd.AddCommand(usersSessionsCmd)
}
