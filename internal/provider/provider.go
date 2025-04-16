package provider

import (
	"errors"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
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

	logger.LogObjectChange(constants.ObjectTypeProvider, constants.ActionCreated, name)

	return provider, nil
}

func DeleteProvider(name string, pk int) (err error) {
	err = ak.Repo.DeleteProvider(pk)
	if err != nil {
		var notExistsError *customErrors.NotExists
		if errors.As(err, &notExistsError) {
			return nil
		}
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeProvider, constants.ActionDeleted, name)

	return nil
}

func GetOidcProvider(pk int) (provider *ak.OidcProvider, err error) {
	provider, err = ak.Repo.GetOidcProvider(pk)
	if err != nil {
		var notExistsError *customErrors.NotExists
		if errors.As(err, &notExistsError) {
			return nil, nil
		}
		return nil, err
	}

	return provider, nil
}
