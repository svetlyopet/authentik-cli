package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/tenant"
)

var CreateTenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "Create a tenant",
	Long: `Tenants are not native objects to Authentik,
but rather an abstraction that we create to bundle applications
and providers for different permissions.

Examples:
  # Create a tenant
  authentik-cli create tenant example-tenant`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tenantName := args[0]
		err := tenant.Create(tenantName)
		cobra.CheckErr(err)
	},
}
