package cmd

import (
	"github.com/spf13/cobra"
	g "github.com/svetlyopet/authentik-cli/cmd/create/group"
	t "github.com/svetlyopet/authentik-cli/cmd/create/tenant"
	u "github.com/svetlyopet/authentik-cli/cmd/create/user"
)

func CreateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a resource from stdin",
		Long: `Creates resources in Authentik which are native,
like applications and providers, or abstractions, like tenants,
which are created by this tool.

Examples:
  # Create a tenant
  authentik-cli create tenant example-tenant`,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			cobra.CheckErr(err)
		},
	}

	c.AddCommand(t.CreateTenantCmd())
	c.AddCommand(u.CreateUserCmd())
	c.AddCommand(g.CreateGroupCmd())

	return c
}
