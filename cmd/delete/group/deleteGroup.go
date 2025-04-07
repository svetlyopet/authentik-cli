package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

func DeleteGroupCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "group",
		Short: "Delete a group",
		Long: `Deletes a group in Authentik.

Examples:
  # Delete a user
  authentik-cli delete group example-group`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			err := core.DeleteGroup(name)
			cobra.CheckErr(err)
		},
	}

	return c
}
