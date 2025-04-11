package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	a "github.com/svetlyopet/authentik-cli/cmd/delete/app"
	g "github.com/svetlyopet/authentik-cli/cmd/delete/group"
	t "github.com/svetlyopet/authentik-cli/cmd/delete/tenant"
	u "github.com/svetlyopet/authentik-cli/cmd/delete/user"
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
			cmd.Help()
		},
	}

	c.AddCommand(t.DeleteTenantCmd())
	c.AddCommand(u.DeleteUserCmd())
	c.AddCommand(g.DeleteGroupCmd())
	c.AddCommand(a.DeleteAppCmd())

	return c
}
