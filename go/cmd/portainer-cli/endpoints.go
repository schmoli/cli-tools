package main

import (
	"github.com/spf13/cobra"
	"github.com/toli/portainer-cli/pkg/portainer"
)

var endpointsCmd = &cobra.Command{
	Use:   "endpoints",
	Short: "Manage endpoints",
}

var endpointsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		endpoints, err := client.ListEndpoints()
		if err != nil {
			handleError(err)
			return
		}

		output := portainer.EndpointList{
			Endpoints: make([]portainer.Endpoint, len(endpoints)),
		}
		for i, e := range endpoints {
			output.Endpoints[i] = e.ToEndpoint()
		}

		if err := portainer.PrintYAML(output); err != nil {
			handleError(err)
		}
	},
}

var endpointsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show an endpoint by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		id, err := parseID(args[0])
		if err != nil {
			handleError(err)
			return
		}

		endpoint, err := client.GetEndpoint(id)
		if err != nil {
			handleError(err)
			return
		}

		if err := portainer.PrintYAML(endpoint.ToEndpoint()); err != nil {
			handleError(err)
		}
	},
}

func init() {
	endpointsCmd.AddCommand(endpointsListCmd)
	endpointsCmd.AddCommand(endpointsShowCmd)
}
