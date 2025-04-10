package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/app"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

var slug string
var providerType string
var oidcClientType string
var oidcRedirectUris []string

func CreateAppCmd() *cobra.Command {
	var err error

	c := &cobra.Command{
		Use:   "app",
		Short: "Create an application and provider pair",
		Long: fmt.Sprintf(`Creates an application and provider in Authentik.

Examples:
  # Create an oidc application
  %s create app example-app \
--slug=example-app \
--provider-type=oidc \
--oidc-client-type=public \
--oidc-redirect-uri=http://localhost:8000`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			err = createApp(name, providerType)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVar(&slug, "slug", "", "Given name of the user")
	c.Flags().StringVar(&providerType, "provider-type", "", "Surname of the user")
	c.Flags().StringVar(&oidcClientType, "oidc-client-type", "", "Email associated with the user")
	c.Flags().StringArrayVar(&oidcRedirectUris, "oidc-redirect-uri", []string{}, "Grant admin permissions for a Tenant")

	err = c.MarkFlagRequired("provider-type")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	return c
}

func createApp(name, providerType string) error {
	switch providerType {
	case "oidc":
		err := app.CreateOidcApplication(name, slug, providerType, oidcClientType, oidcRedirectUris)
		if err != nil {
			return err
		}
	case "ldap":
		return fmt.Errorf("LDAP provider is not supported yet")
	case "saml":
		return fmt.Errorf("SAML provider is not supported yet")
	case "proxy":
		return fmt.Errorf("Proxy provider is not supported yet")
	default:
		return fmt.Errorf("%s is not a supported provider type", providerType)
	}

	return nil
}
