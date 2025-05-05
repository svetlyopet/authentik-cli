package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/logger"
	"github.com/svetlyopet/authentik-cli/internal/tenant"
	"gopkg.in/yaml.v3"
)

func GetTenantCmd() *cobra.Command {
	var err error
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
			err = GetTenantWithFormat(name, outputFormat)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json|yaml)")
	if outputFormat != "json" && outputFormat != "yaml" {
		cobra.CheckErr(fmt.Errorf("invalid output format: %s. Supported formats: json, yaml", outputFormat))
	}

	return c
}

func GetTenantWithFormat(name, outputFormat string) error {
	var tenantDetailsMarshaled []byte

	tenantDetails, err := tenant.Get(name)
	if err != nil {
		return err
	}

	if outputFormat == "json" {
		tenantDetailsMarshaled, err = json.MarshalIndent(tenantDetails, "", "  ")
		if err != nil {
			return err
		}
	}

	if outputFormat == "yaml" {
		tenantDetailsMarshaled, err = yaml.Marshal(tenantDetails)
		if err != nil {
			return err
		}
	}

	logger.LogObjectDetails(tenantDetailsMarshaled)

	return nil
}
