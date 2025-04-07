package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/tenant"
)

func CreateTenantCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "tenant",
		Short: "Create a tenant",
		Long: fmt.Sprintf(`Tenants are not native objects to Authentik,
but rather an abstraction that we create to bundle applications
and providers for different permissions.

Examples:
  # Create a tenant
  %s create tenant example-tenant`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tenantName := args[0]
			err := tenant.Create(tenantName)
			cobra.CheckErr(err)
		},
	}

	return c
}
