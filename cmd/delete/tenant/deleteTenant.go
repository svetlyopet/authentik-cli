package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/tenant"
)

var DeleteTenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "Delete a tenant",
	Long: `Tenants are not native objects to Authentik,
but rather an abstraction that we create to bundle applications
and providers for different permissions.

Examples:
  # Delete a tenant
  authentik-cli delete tenant example-tenant`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tenantName := args[0]
		err := tenant.Delete(tenantName)
		cobra.CheckErr(err)
	},
}
