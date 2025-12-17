// keycloak/cmd/keycloak-cli/roles.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var (
	flagClientUUID string
)

var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Manage roles",
}

var rolesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List roles (realm or client)",
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

		var roles *keycloak.RoleList
		if flagClientUUID != "" {
			roles, err = client.ListClientRoles(realm, flagClientUUID)
		} else {
			roles, err = client.ListRealmRoles(realm)
		}
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(roles); err != nil {
			handleError(err)
		}
	},
}

var rolesGetCmd = &cobra.Command{
	Use:   "get <role-name>",
	Short: "Get role details",
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

		var role *keycloak.RoleInfo
		if flagClientUUID != "" {
			role, err = client.GetClientRole(realm, flagClientUUID, args[0])
		} else {
			role, err = client.GetRealmRole(realm, args[0])
		}
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(role); err != nil {
			handleError(err)
		}
	},
}

func init() {
	rolesCmd.PersistentFlags().StringVar(&flagClientUUID, "client", "", "Client UUID for client roles")
	rolesCmd.AddCommand(rolesListCmd)
	rolesCmd.AddCommand(rolesGetCmd)
}
