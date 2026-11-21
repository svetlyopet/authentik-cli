package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/app"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/provider"
)

func updateAppCmd() *cobra.Command {
	var name string
	var slug string
	var oidcClientType string
	var oidcConsentType string
	var oidcRedirectUrisStrict []string
	var oidcRedirectUrisRegex []string

	c := &cobra.Command{
		Use:   "app",
		Short: "Update an application and its provider",
		Long: fmt.Sprintf(`Updates an application and its provider in Authentik.
Only the attributes passed as flags are changed.

Renaming an application also renames its provider when the provider still
carries the application's name. Changing the slug changes the application's
OIDC endpoints and, with the per-provider issuer mode, its issuer URL.

Each --oidc-redirect-uri-* flag replaces only the redirect URIs of its own
matching mode; pass it with an empty value to remove them all.

Examples:
  # Make an OIDC application public
  %[1]s update app example-app --oidc-client-type=public

  # Replace the strict redirect URIs, keeping the regex ones
  %[1]s update app example-app \
--oidc-redirect-uri-strict=https://app.example.com/callback \
--oidc-redirect-uri-strict=http://localhost:8000/callback

  # Remove all regex redirect URIs
  %[1]s update app example-app --oidc-redirect-uri-regex=

  # Rename an application
  %[1]s update app example-app --name=new-app-name`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			changes := app.Changes{
				Name: flagValue(cmd, "name", name),
				Slug: flagValue(cmd, "slug", slug),
				Oidc: provider.OidcChanges{
					ClientType:         flagValue(cmd, "oidc-client-type", oidcClientType),
					ConsentType:        flagValue(cmd, "oidc-consent-type", oidcConsentType),
					RedirectUrisStrict: flagValue(cmd, "oidc-redirect-uri-strict", oidcRedirectUrisStrict),
					RedirectUrisRegex:  flagValue(cmd, "oidc-redirect-uri-regex", oidcRedirectUrisRegex),
				},
			}
			err := app.Update(args[0], changes)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVar(&name, "name", "", "New name of the application")
	c.Flags().StringVar(&slug, "slug", "", "New slug of the application")
	c.Flags().StringVar(&oidcClientType, "oidc-client-type", "", "OIDC client type (confidential|public)")
	c.Flags().StringVar(&oidcConsentType, "oidc-consent-type", "", "OIDC consent type (explicit|implicit)")
	c.Flags().StringArrayVar(&oidcRedirectUrisStrict, "oidc-redirect-uri-strict", []string{}, "Strict redirect URIs for the OIDC provider (replaces the current strict URIs)")
	c.Flags().StringArrayVar(&oidcRedirectUrisRegex, "oidc-redirect-uri-regex", []string{}, "Regex redirect URIs for the OIDC provider (replaces the current regex URIs)")

	c.MarkFlagsOneRequired("name", "slug", "oidc-client-type", "oidc-consent-type",
		"oidc-redirect-uri-strict", "oidc-redirect-uri-regex")

	return c
}
