package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

func updateGroupCmd() *cobra.Command {
	var name string

	c := &cobra.Command{
		Use:   "group",
		Short: "Update a group",
		Long: fmt.Sprintf(`Updates a local group in Authentik. Only the attributes
passed as flags are changed. Admin groups of tenants cannot be renamed.

Examples:
  # Rename a group
  %s update group example-group --name=new-group-name`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			changes := core.GroupChanges{
				Name: flagValue(cmd, "name", name),
			}
			err := core.UpdateGroup(args[0], changes)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVar(&name, "name", "", "New name of the group")

	c.MarkFlagsOneRequired("name")

	return c
}
