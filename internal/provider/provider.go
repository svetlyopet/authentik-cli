package provider

import (
	"errors"
	"fmt"
	"slices"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/diff"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	"github.com/svetlyopet/authentik-cli/internal/flow"
	"github.com/svetlyopet/authentik-cli/internal/logger"
)

func CreateOidcProvider(name, oidcClientType, oidcConsentType string, oidcEncryptTokens bool, oidcRedirectUrisStrict, oidcRedirectUrisRegex []string) (provider *ak.OidcProvider, err error) {
	var redirectUris []ak.OidcRedirectUri
	redirectUris = append(redirectUris, toRedirectUris(OidcRedirectUriMatchingModeStrict, oidcRedirectUrisStrict)...)
	redirectUris = append(redirectUris, toRedirectUris(OidcRedirectUriMatchingModeRegex, oidcRedirectUrisRegex)...)

	authorizationFlowSlug, err := authorizationFlowSlugFor(oidcConsentType)
	if err != nil {
		return nil, err
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

// OidcChanges lists the changes requested for an OIDC provider. A nil field
// is left unchanged.
type OidcChanges struct {
	Name        *string
	ClientType  *string
	ConsentType *string
	// RedirectUrisStrict replaces the provider's strict redirect URIs and
	// RedirectUrisRegex its regex ones; the other kind is kept. Empty strings
	// are ignored, so a list holding only "" clears that kind.
	RedirectUrisStrict *[]string
	RedirectUrisRegex  *[]string
}

// IsEmpty reports whether no change is requested.
func (c OidcChanges) IsEmpty() bool {
	return c.Name == nil && c.ClientType == nil && c.ConsentType == nil &&
		c.RedirectUrisStrict == nil && c.RedirectUrisRegex == nil
}

// UpdateOidcProvider applies the requested changes to the OIDC provider with
// the given primary key. Fields that were not requested, or already hold the
// requested value, are not sent to authentik and keep their current value.
func UpdateOidcProvider(pk int, changes OidcChanges) error {
	if changes.ClientType != nil {
		if err := validateClientType(*changes.ClientType); err != nil {
			return err
		}
	}

	var authorizationFlow *string
	if changes.ConsentType != nil {
		slug, err := authorizationFlowSlugFor(*changes.ConsentType)
		if err != nil {
			return err
		}
		flowID, err := flowPK(slug, "authorization", true)
		if err != nil {
			return err
		}
		authorizationFlow = &flowID
	}

	current, err := ak.Repo.GetOidcProvider(pk)
	if err != nil {
		return err
	}

	patch := ak.OidcProviderPatch{
		Name:              diff.Changed(current.Name, changes.Name),
		ClientType:        diff.Changed(current.ClientType, changes.ClientType),
		AuthorizationFlow: diff.Changed(current.AuthorizationFlow, authorizationFlow),
		RedirectUris:      redirectUrisPatch(current.RedirectUris, changes.RedirectUrisStrict, changes.RedirectUrisRegex),
	}
	if patch.IsEmpty() {
		logger.LogObjectChange(constants.ObjectTypeProvider, constants.ActionUnchanged, current.Name)
		return nil
	}

	updated, err := ak.Repo.UpdateOidcProvider(pk, patch)
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeProvider, constants.ActionChanged, updated.Name)

	return nil
}

// redirectUrisPatch returns the current redirect URIs with the strict and/or
// regex ones replaced by the requested URLs, or nil when neither kind is
// requested or the result equals the current list.
func redirectUrisPatch(current []ak.OidcRedirectUri, strict, regex *[]string) *[]ak.OidcRedirectUri {
	if strict == nil && regex == nil {
		return nil
	}

	pick := func(mode string, requested *[]string) []ak.OidcRedirectUri {
		if requested != nil {
			return toRedirectUris(mode, slices.DeleteFunc(slices.Clone(*requested), func(url string) bool {
				return url == ""
			}))
		}
		return slices.DeleteFunc(slices.Clone(current), func(uri ak.OidcRedirectUri) bool {
			return uri.MatchingMode != mode
		})
	}

	uris := slices.Concat(
		pick(OidcRedirectUriMatchingModeStrict, strict),
		pick(OidcRedirectUriMatchingModeRegex, regex),
	)
	if slices.Equal(uris, current) {
		return nil
	}
	if uris == nil {
		// An explicit empty list clears the URIs; nil would mean "unchanged".
		uris = []ak.OidcRedirectUri{}
	}

	return &uris
}

func toRedirectUris(mode string, urls []string) []ak.OidcRedirectUri {
	var uris []ak.OidcRedirectUri
	for _, url := range urls {
		uris = append(uris, ak.OidcRedirectUri{MatchingMode: mode, Url: url})
	}
	return uris
}

// authorizationFlowSlugFor returns the slug of the default authorization flow
// matching an OIDC consent type.
func authorizationFlowSlugFor(consentType string) (string, error) {
	switch consentType {
	case OidcConsentTypeExplicit:
		return FlowSlugAuthorizationExplicitConsent, nil
	case OidcConsentTypeImplicit:
		return FlowSlugAuthorizationImplicitConsent, nil
	default:
		return "", fmt.Errorf("invalid OIDC consent type %q: must be %q or %q",
			consentType, OidcConsentTypeExplicit, OidcConsentTypeImplicit)
	}
}

func validateClientType(clientType string) error {
	switch clientType {
	case OidcClientTypeConfidential, OidcClientTypePublic:
		return nil
	default:
		return fmt.Errorf("invalid OIDC client type %q: must be %q or %q",
			clientType, OidcClientTypeConfidential, OidcClientTypePublic)
	}
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
