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
var oidcEncryptToken bool
var oidcConsentType string
var oidcRedirectUrisStrict []string
var oidcRedirectUrisRegex []string

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
	c.Flags().StringArrayVar(&oidcRedirectUrisStrict, "oidc-redirect-uri-strict", []string{}, "Strict redirect URIs for the OIDC provider")
	c.Flags().StringArrayVar(&oidcRedirectUrisRegex, "oidc-redirect-uri-regex", []string{}, "Regex redirect URIs for the OIDC provider")

	err = c.MarkFlagRequired("provider-type")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	return c
}

func createApp(name, providerType string) (err error) {
	switch providerType {
	case constants.ProviderTypeOIDC:
		err = app.CreateOidc(name, slug, oidcClientType, oidcConsentType, oidcEncryptToken, oidcRedirectUrisStrict, oidcRedirectUrisRegex)
		if err != nil {
			return err
		}
	case constants.ProviderTypeLDAP:
		return fmt.Errorf("LDAP provider is not supported yet")
	case constants.ProviderTypeSAML:
		return fmt.Errorf("SAML provider is not supported yet")
	case constants.ProviderTypeProxy:
		return fmt.Errorf("Proxy provider is not supported yet")
	case constants.ProviderTypeRAC:
		return fmt.Errorf("RAC provider is not supported yet")
	case constants.ProviderTypeSCIM:
		return fmt.Errorf("SCIM provider is not supported yet")
	default:
		return fmt.Errorf("%s is not a supported provider type", providerType)
	}

	return nil
}
