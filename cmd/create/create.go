package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	a "github.com/svetlyopet/authentik-cli/cmd/create/app"
	g "github.com/svetlyopet/authentik-cli/cmd/create/group"
	t "github.com/svetlyopet/authentik-cli/cmd/create/tenant"
	u "github.com/svetlyopet/authentik-cli/cmd/create/user"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

func CreateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a resource from stdin",
		Long: fmt.Sprintf(`Creates resources in Authentik which are native,
like applications and providers, or abstractions, like tenants,
which are created by this tool.

Examples:
  # Create a tenant
  %s create tenant example-tenant`, constants.CmdName),
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			cobra.CheckErr(err)
		},
	}

	c.AddCommand(t.CreateTenantCmd())
	c.AddCommand(u.CreateUserCmd())
	c.AddCommand(g.CreateGroupCmd())
	c.AddCommand(a.CreateAppCmd())

	return c
}
