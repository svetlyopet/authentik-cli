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

	redirectUris := make([]api.RedirectURIRequest, 0, len(oidcProvider.RedirectUris))
	for _, uri := range oidcProvider.RedirectUris {
		redirectUris = append(redirectUris, *api.NewRedirectURIRequest(api.MatchingModeEnum(uri.MatchingMode), uri.Url))
	}

	req := api.NewOAuth2ProviderRequest(
		oidcProvider.Name,
		oidcProvider.AuthorizationFlow,
		oidcProvider.InvalidationFlow,
		redirectUris,
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
