package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/tenant"
)

func getTenantCmd() *cobra.Command {
	var outputFormat string

	c := &cobra.Command{
		Use:   "tenant",
		Short: "Get details about a tenant",
		Long: fmt.Sprintf(`Tenants are not native objects to Authentik,
but rather an abstraction that we create to bundle applications
and providers for different permissions.

Examples:
  # Get details about a tenant
  %s get tenant example-tenant`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			tenantDetails, err := tenant.Get(name)
			cobra.CheckErr(err)
			outputDetailsWithFormat(tenantDetails, outputFormat)
		},
	}

	c.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json|yaml)")
	if outputFormat != "json" && outputFormat != "yaml" {
		cobra.CheckErr(fmt.Errorf("invalid output format: %s. Supported formats: json, yaml", outputFormat))
	}

	return c
}
