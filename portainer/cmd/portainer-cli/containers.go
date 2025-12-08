package main

import (
	"sort"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/portainer/pkg/portainer"
)

var flagEndpoint int64

var containersCmd = &cobra.Command{
	Use:   "containers",
	Short: "Manage containers",
}

var containersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all containers",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		var endpointIDs []int64

		if flagEndpoint > 0 {
			// Use specified endpoint
			endpointIDs = []int64{flagEndpoint}
		} else {
			// Get all endpoints
			endpoints, err := client.ListEndpoints()
			if err != nil {
				handleError(err)
				return
			}
			for _, e := range endpoints {
				endpointIDs = append(endpointIDs, e.ID)
			}
		}

		var output portainer.ContainerList

		for _, eid := range endpointIDs {
			containers, err := client.ListContainers(eid)
			if err != nil {
				handleError(err)
				return
			}
			for _, c := range containers {
				output = append(output, c.ToListItem(eid))
			}
		}

		// Sort by endpoint, then name
		sort.Slice(output, func(i, j int) bool {
			if output[i].Endpoint != output[j].Endpoint {
				return output[i].Endpoint < output[j].Endpoint
			}
			return output[i].Name < output[j].Name
		})

		if err := portainer.PrintYAML(output); err != nil {
			handleError(err)
		}
	},
}

func init() {
	containersListCmd.Flags().Int64Var(&flagEndpoint, "endpoint", 0, "Filter by endpoint ID")
	containersCmd.AddCommand(containersListCmd)
}
