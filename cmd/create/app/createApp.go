package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
	"github.com/svetlyopet/authentik-cli/internal/provider"
)

var slug string
var providerType string
var oidcClientType string
var oidcEncryptToken bool
var oidcConsentType string
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

	c.Flags().StringVar(&slug, "slug", "", "Application slug")
	c.Flags().StringVar(&providerType, "provider-type", "", "Provider type")
	c.Flags().StringVar(&oidcClientType, "oidc-client-type", "public", "OIDC client type")
	c.Flags().StringVar(&oidcConsentType, "oidc-consent-type", "explicit", "OIDC consent type")
	c.Flags().BoolVar(&oidcEncryptToken, "oidc-encrypt-tokens", false, "Enable encrypted tokens for OIDC provider")
	c.Flags().StringArrayVar(&oidcRedirectUris, "oidc-redirect-uri", []string{}, "Redirect URIs for the OIDC provider")

	err = c.MarkFlagRequired("provider-type")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	return c
}

func createApp(name, providerType string) error {
	var providerPK int

	switch providerType {
	case "oidc":
		p, err := provider.CreateOidcProvider(name, oidcClientType, oidcConsentType, oidcEncryptToken, oidcRedirectUris)
		if err != nil {
			return err
		}
		providerPK = p.PK
	case "ldap":
		return fmt.Errorf("LDAP provider is not supported yet")
	case "saml":
		return fmt.Errorf("SAML provider is not supported yet")
	case "proxy":
		return fmt.Errorf("Proxy provider is not supported yet")
	default:
		return fmt.Errorf("%s is not a supported provider type", providerType)
	}

	err := core.CreateApplication(name, slug, providerPK)
	if err != nil {
		return err
	}

	return nil
}
