// keycloak/cmd/keycloak-cli/realms.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var realmsCmd = &cobra.Command{
	Use:   "realms",
	Short: "Manage realms",
}

var realmsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all realms",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		realms, err := client.ListRealms()
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(realms); err != nil {
			handleError(err)
		}
	},
}

var realmsGetCmd = &cobra.Command{
	Use:   "get <realm-name>",
	Short: "Get realm details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		realm, err := client.GetRealm(args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(realm); err != nil {
			handleError(err)
		}
	},
}

func init() {
	realmsCmd.AddCommand(realmsListCmd)
	realmsCmd.AddCommand(realmsGetCmd)
}
