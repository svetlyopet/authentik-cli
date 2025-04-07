package cmd

import (
	"github.com/spf13/cobra"
	t "github.com/svetlyopet/authentik-cli/cmd/delete/tenant"
	u "github.com/svetlyopet/authentik-cli/cmd/delete/user"
)

func DeleteCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource from stdin",
		Long: `Deletes resources in Authentik which are native,
like applications and providers, or abstractions, like tenants,
which are created by this tool.

Examples:
  # Delete a tenant
  authentik-cli delete tenant example-tenant`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	c.AddCommand(t.DeleteTenantCmd())
	c.AddCommand(u.DeleteUserCmd())

	return c
}
