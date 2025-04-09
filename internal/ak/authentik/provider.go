package authentik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
)

const (
	providersOAuth2Path = "%s/api/v3/providers/oauth2/"
)

func (a *authentik) CreateOidcProvider(oidcProvider ak.OidcProvider) (*ak.OidcProvider, error) {
	var redirectUris []oidcRedirectUri
	for _, redirectUri := range oidcProvider.RedirectUris {
		redirectUris = append(redirectUris, oidcRedirectUri{
			MatchingMode: redirectUri.MatchingMode,
			Url:          redirectUri.Url,
		})
	}

	createOrUpdateOidcProviderReq := createOrUpdateOidcProviderRequest{
		Name:                 oidcProvider.Name,
		AuthenticationFlow:   oidcProvider.AuthenticationFlow,
		AuthorizationFlow:    oidcProvider.AuthorizationFlow,
		InvalidationFlow:     oidcProvider.InvalidationFlow,
		PropertyMappings:     oidcProvider.PropertyMappings,
		ClientType:           oidcProvider.ClientType,
		AccessCodeValidity:   oidcProvider.AccessCodeValidity,
		AccessTokenValidity:  oidcProvider.AccessTokenValidity,
		RefreshTokenValidity: oidcProvider.RefreshTokenValidity,
		SigningKey:           oidcProvider.SigningKey,
		EncryptionKey:        oidcProvider.EncryptionKey,
		RedirectUris:         redirectUris,
		SubMode:              oidcProvider.SubMode,
		IssuerMode:           oidcProvider.IssuerMode,
	}

	createOrUpdateOidcProviderRequestBytes, err := json.Marshal(createOrUpdateOidcProviderReq)
	if err != nil {
		return nil, err
	}

	response, err := a.doRequest(http.MethodPost,
		fmt.Sprintf(providersOAuth2Path, a.url),
		bytes.NewBuffer(createOrUpdateOidcProviderRequestBytes))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(response.Body)
		return nil, customErrors.NewUnexpectedResult(fmt.Sprintf("create oidc provider: %s", string(errBody)))
	}

	var oidcProvidersResp createOrUpdateOidcProviderResponse
	err = json.NewDecoder(response.Body).Decode(&oidcProvidersResp)
	if err != nil {
		return nil, err
	}

	return mapToCreateOrUpdateProviderResponse(&oidcProvidersResp), nil
}
