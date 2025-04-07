package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

func DeleteUserCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "user",
		Short: "Delete a user",
		Long: `Deletes a user in Authentik.

Examples:
  # Delete a user
  authentik-cli delete user example-user`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			err := core.DeleteUser(username)
			cobra.CheckErr(err)
		},
	}

	return c
}
