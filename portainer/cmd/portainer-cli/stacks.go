package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/portainer/pkg/portainer"
)

var stacksCmd = &cobra.Command{
	Use:   "stacks",
	Short: "Manage stacks",
}

var stacksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stacks",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		stacks, err := client.ListStacks()
		if err != nil {
			handleError(err)
			return
		}

		output := portainer.StackList{
			Stacks: make([]portainer.StackListItem, len(stacks)),
		}
		for i, s := range stacks {
			output.Stacks[i] = s.ToListItem()
		}

		if err := portainer.PrintYAML(output); err != nil {
			handleError(err)
		}
	},
}

var stacksShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a stack by ID",
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

		// Get stack list to find endpoint_id
		stacks, err := client.ListStacks()
		if err != nil {
			handleError(err)
			return
		}

		var apiStack *portainer.APIStack
		for i := range stacks {
			if stacks[i].ID == id {
				apiStack = &stacks[i]
				break
			}
		}
		if apiStack == nil {
			handleError(portainer.NotFoundError(fmt.Sprintf("stack with ID %d", id)))
			return
		}

		// Get stack file content
		file, err := client.GetStackFile(id)
		if err != nil {
			handleError(err)
			return
		}

		stack := apiStack.ToStack(file.StackFileContent)
		if err := portainer.PrintYAML(stack); err != nil {
			handleError(err)
		}
	},
}

func init() {
	stacksCmd.AddCommand(stacksListCmd)
	stacksCmd.AddCommand(stacksShowCmd)
}
