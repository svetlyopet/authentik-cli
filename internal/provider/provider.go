package provider

import (
	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/flow"
	"github.com/svetlyopet/authentik-cli/internal/logger"
)

func CreateOidcProvider(name, oidcClientType, oidcConsentType string, oidcEncryptTokens bool, oidcRedirectUris []string) (provider *ak.OidcProvider, err error) {
	// TODO: add support for regex urls
	var redirectUris []ak.OidcRedirectUri
	for _, url := range oidcRedirectUris {
		redirectUris = append(redirectUris, ak.OidcRedirectUri{
			MatchingMode: OidcRedirectUriMatchingModeStrict,
			Url:          url,
		})
	}

	var authenticationFlow string
	var authorizationFlow string
	var invalidationFlow string

	flows, err := flow.GetFlows()
	if err != nil {
		return nil, err
	}

	// TODO: add options for choosing the flows instead of always using the defaults
	for _, f := range flows {
		if f.Designation == "authorization" {
			if oidcConsentType == "explicit" && f.Slug == "default-provider-authorization-explicit-consent" {
				authorizationFlow = f.PK
			} else if oidcConsentType == "implicit" && f.Slug == "default-provider-authorization-implicit-consent" {
				authorizationFlow = f.PK
			}
		} else if f.Designation == "authentication" && f.Slug == "default-authentication-flow" {
			authenticationFlow = f.PK
		} else if f.Designation == "invalidation" && f.Slug == "default-invalidation-flow" {
			invalidationFlow = f.PK
		}
	}

	// TODO: add support for signing and encryption of tokens, and property mappings
	oidcProvider := ak.OidcProvider{
		Provider: ak.Provider{
			Name:               name,
			AuthenticationFlow: authenticationFlow,
			AuthorizationFlow:  authorizationFlow,
			InvalidationFlow:   invalidationFlow,
		},
		ClientType:           oidcClientType,
		AccessCodeValidity:   OidcAccessCodeValidityDefault,
		AccessTokenValidity:  OidcAccessTokenValidityDefault,
		RefreshTokenValidity: OidcRefreshTokenValidityDefault,
		RedirectUris:         redirectUris,
		SubMode:              OidcSubModeDefault,
		IssuerMode:           OidcIssuerModeDefault,
	}

	provider, err = ak.Repo.CreateOidcProvider(oidcProvider)
	if err != nil {
		return nil, err
	}

	logger.WriteStdout(constants.ObjectTypeProvider, constants.ActionCreated, name)

	return provider, nil
}
