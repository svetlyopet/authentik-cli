package provider

import (
	"errors"
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	"github.com/svetlyopet/authentik-cli/internal/flow"
	"github.com/svetlyopet/authentik-cli/internal/logger"
)

func CreateOidcProvider(name, oidcClientType, oidcConsentType string, oidcEncryptTokens bool, oidcRedirectUrisStrict, oidcRedirectUrisRegex []string) (provider *ak.OidcProvider, err error) {
	var redirectUrisStrict []ak.OidcRedirectUri
	for _, url := range oidcRedirectUrisStrict {
		redirectUrisStrict = append(redirectUrisStrict, ak.OidcRedirectUri{
			MatchingMode: OidcRedirectUriMatchingModeStrict,
			Url:          url,
		})
	}
	var redirectUrisRegex []ak.OidcRedirectUri
	for _, url := range oidcRedirectUrisRegex {
		redirectUrisRegex = append(redirectUrisRegex, ak.OidcRedirectUri{
			MatchingMode: OidcRedirectUriMatchingModeRegex,
			Url:          url,
		})
	}

	var redirectUris []ak.OidcRedirectUri
	redirectUris = append(redirectUris, redirectUrisStrict...)
	redirectUris = append(redirectUris, redirectUrisRegex...)

	var authorizationFlowSlug string
	switch oidcConsentType {
	case OidcConsentTypeExplicit:
		authorizationFlowSlug = FlowSlugAuthorizationExplicitConsent
	case OidcConsentTypeImplicit:
		authorizationFlowSlug = FlowSlugAuthorizationImplicitConsent
	default:
		return nil, fmt.Errorf("invalid OIDC consent type %q: must be %q or %q",
			oidcConsentType, OidcConsentTypeExplicit, OidcConsentTypeImplicit)
	}

	// TODO: add options for choosing the flows instead of always using the defaults
	authenticationFlow, err := flowPK(FlowSlugAuthentication, "authentication", false)
	if err != nil {
		return nil, err
	}
	authorizationFlow, err := flowPK(authorizationFlowSlug, "authorization", true)
	if err != nil {
		return nil, err
	}
	invalidationFlow, err := flowPK(FlowSlugInvalidation, "invalidation", true)
	if err != nil {
		return nil, err
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

// flowPK looks up the flow with the given slug and returns its primary key.
// A missing optional flow yields an empty key; a missing required flow, or a
// flow with an unexpected designation, is an error.
func flowPK(slug, designation string, required bool) (string, error) {
	f, err := flow.GetFlowBySlug(slug)
	if err != nil {
		var notExistsError *customErrors.NotExists
		if errors.As(err, &notExistsError) {
			if !required {
				return "", nil
			}
			return "", customErrors.NewNotExists(fmt.Sprintf(
				"required flow %q not found in authentik (make sure the default blueprints have been applied)", slug))
		}
		return "", err
	}

	if f.Designation != designation {
		return "", fmt.Errorf("flow %q has designation %q, expected %q", slug, f.Designation, designation)
	}
	if f.PK == "" {
		return "", fmt.Errorf("flow %q was returned without a primary key", slug)
	}

	return f.PK, nil
}
