// keycloak/cmd/keycloak-cli/groups.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage groups",
}

var groupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List groups in realm",
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

		groups, err := client.ListGroups(realm)
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(groups); err != nil {
			handleError(err)
		}
	},
}

var groupsGetCmd = &cobra.Command{
	Use:   "get <group-id>",
	Short: "Get group details",
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

		group, err := client.GetGroup(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(group); err != nil {
			handleError(err)
		}
	},
}

var groupsMembersCmd = &cobra.Command{
	Use:   "members <group-id>",
	Short: "List group members",
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

		members, err := client.GetGroupMembers(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(members); err != nil {
			handleError(err)
		}
	},
}

func init() {
	groupsCmd.AddCommand(groupsListCmd)
	groupsCmd.AddCommand(groupsGetCmd)
	groupsCmd.AddCommand(groupsMembersCmd)
}
