package cmd

import (
	"github.com/spf13/cobra"
	t "github.com/svetlyopet/authentik-cli/cmd/create/tenant"
	u "github.com/svetlyopet/authentik-cli/cmd/create/user"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from stdin",
	Long: `Creates resources in Authentik which are native,
like applications and providers, or abstractions, like tenants,
which are created by this tool.

Examples:
  # Create a tenant
  authentik-cli create tenant example-tenant`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func addSubcommands() {
	CreateCmd.AddCommand(t.CreateTenantCmd)
	CreateCmd.AddCommand(u.CreateUserCmd)
}

func init() {
	addSubcommands()
}
