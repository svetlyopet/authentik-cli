package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

func DeleteCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource from stdin",
		Long: fmt.Sprintf(`Deletes resources in Authentik which are native,
like applications and providers, or abstractions, like tenants,
which are created by this tool.

Examples:
  # Delete a tenant
  %s delete tenant example-tenant`, constants.CmdName),
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			cobra.CheckErr(err)
		},
	}

	c.AddCommand(deleteTenantCmd())
	c.AddCommand(deleteUserCmd())
	c.AddCommand(deleteGroupCmd())
	c.AddCommand(deleteAppCmd())

	return c
}
