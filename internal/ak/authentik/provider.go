package authentik

import (
	"context"

	api "goauthentik.io/api/v3"

	"github.com/svetlyopet/authentik-cli/internal/ak"
)

func (a *authentik) CreateOidcProvider(oidcProvider ak.OidcProvider) (*ak.OidcProvider, error) {
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	req := api.NewOAuth2ProviderRequest(
		oidcProvider.Name,
		oidcProvider.AuthorizationFlow,
		oidcProvider.InvalidationFlow,
		redirectUrisToAPI(oidcProvider.RedirectUris),
	)
	if oidcProvider.AuthenticationFlow != "" {
		req.SetAuthenticationFlow(oidcProvider.AuthenticationFlow)
	}
	if len(oidcProvider.PropertyMappings) > 0 {
		req.SetPropertyMappings(oidcProvider.PropertyMappings)
	}
	if oidcProvider.ClientType != "" {
		req.SetClientType(api.ClientTypeEnum(oidcProvider.ClientType))
	}
	if oidcProvider.AccessCodeValidity != "" {
		req.SetAccessCodeValidity(oidcProvider.AccessCodeValidity)
	}
	if oidcProvider.AccessTokenValidity != "" {
		req.SetAccessTokenValidity(oidcProvider.AccessTokenValidity)
	}
	if oidcProvider.RefreshTokenValidity != "" {
		req.SetRefreshTokenValidity(oidcProvider.RefreshTokenValidity)
	}
	if oidcProvider.SigningKey != "" {
		req.SetSigningKey(oidcProvider.SigningKey)
	}
	if oidcProvider.EncryptionKey != "" {
		req.SetEncryptionKey(oidcProvider.EncryptionKey)
	}
	if oidcProvider.SubMode != "" {
		req.SetSubMode(api.SubModeEnum(oidcProvider.SubMode))
	}
	if oidcProvider.IssuerMode != "" {
		req.SetIssuerMode(api.IssuerModeEnum(oidcProvider.IssuerMode))
	}

	provider, resp, err := a.client.ProvidersAPI.ProvidersOauth2Create(context.Background()).
		OAuth2ProviderRequest(*req).
		Execute()
	if err != nil {
		return nil, apiError("create oidc provider", resp, err)
	}

	return mapOidcProvider(provider), nil
}

func (a *authentik) DeleteProvider(pk int) error {
	const op = "delete provider"
	if err := a.preflightCheck(); err != nil {
		return err
	}

	id, err := toInt32(op, pk)
	if err != nil {
		return err
	}

	resp, err := a.client.ProvidersAPI.ProvidersAllDestroy(context.Background(), id).Execute()
	if err != nil {
		return apiError(op, resp, err)
	}

	return nil
}

func (a *authentik) GetOidcProvider(pk int) (*ak.OidcProvider, error) {
	const op = "get oidc provider"
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	id, err := toInt32(op, pk)
	if err != nil {
		return nil, err
	}

	provider, resp, err := a.client.ProvidersAPI.ProvidersOauth2Retrieve(context.Background(), id).Execute()
	if err != nil {
		return nil, apiError(op, resp, err)
	}

	return mapOidcProvider(provider), nil
}

// UpdateOidcProvider applies a partial update to the OAuth2/OIDC provider with
// the given primary key. Only the fields set in patch are sent; authentik
// keeps every other field.
func (a *authentik) UpdateOidcProvider(pk int, patch ak.OidcProviderPatch) (*ak.OidcProvider, error) {
	const op = "update oidc provider"
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	id, err := toInt32(op, pk)
	if err != nil {
		return nil, err
	}

	// Start from the zero value: the generated New* constructors fill in schema
	// defaults (client type, sub mode, token validity, ...), which a PATCH
	// would then write over the provider's current settings.
	req := api.PatchedOAuth2ProviderRequest{}
	if patch.Name != nil {
		req.SetName(*patch.Name)
	}
	if patch.ClientType != nil {
		req.SetClientType(api.ClientTypeEnum(*patch.ClientType))
	}
	if patch.AuthorizationFlow != nil {
		req.SetAuthorizationFlow(*patch.AuthorizationFlow)
	}
	if patch.RedirectUris != nil {
		req.SetRedirectUris(redirectUrisToAPI(*patch.RedirectUris))
	}

	provider, resp, err := a.client.ProvidersAPI.ProvidersOauth2PartialUpdate(context.Background(), id).
		PatchedOAuth2ProviderRequest(req).
		Execute()
	if err != nil {
		return nil, apiError(op, resp, err)
	}

	return mapOidcProvider(provider), nil
}

// redirectUrisToAPI converts redirect URIs to their API form. The result is
// never nil, so an empty list is still sent and clears the provider's URIs.
func redirectUrisToAPI(uris []ak.OidcRedirectUri) []api.RedirectURIRequest {
	out := make([]api.RedirectURIRequest, 0, len(uris))
	for _, uri := range uris {
		out = append(out, *api.NewRedirectURIRequest(api.MatchingModeEnum(uri.MatchingMode), uri.Url))
	}
	return out
}
