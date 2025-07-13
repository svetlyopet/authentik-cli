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
	providersOAuth2Path             = "%s/api/v3/providers/oauth2/"
	providersOAuth2UpdateDeletePath = "%s/api/v3/providers/oauth2/%d/"
	providersAllPath                = "%s/api/v3/providers/all/"
	providersAllDeletePath          = "%s/api/v3/providers/all/%d/"
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
	defer response.Body.Close() //nolint

	if response.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(response.Body)
		return nil, customErrors.NewUnexpectedResult(fmt.Sprintf("create oidc provider: %s", string(errBody)))
	}

	var oidcProvidersResp getOidcProviderResponse
	err = json.NewDecoder(response.Body).Decode(&oidcProvidersResp)
	if err != nil {
		return nil, err
	}

	return mapToGetOidcProviderResponse(&oidcProvidersResp), nil
}

func (a *authentik) DeleteProvider(id int) error {
	response, err := a.doRequest(http.MethodDelete, fmt.Sprintf(providersAllDeletePath, a.url, id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close() //nolint

	if response.StatusCode == http.StatusNotFound {
		return customErrors.NewNotExists("provider not found")
	}

	if response.StatusCode != http.StatusNoContent {
		errBody, _ := io.ReadAll(response.Body)
		return customErrors.NewUnexpectedResult(fmt.Sprintf("delete provider: %s", string(errBody)))
	}

	return nil
}

func (a *authentik) GetOidcProvider(id int) (*ak.OidcProvider, error) {
	response, err := a.doRequest(http.MethodGet, fmt.Sprintf(providersOAuth2UpdateDeletePath, a.url, id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() //nolint

	if response.StatusCode == http.StatusNotFound {
		return nil, customErrors.NewNotExists("provider not found")
	}

	if response.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(response.Body)
		return nil, customErrors.NewUnexpectedResult(fmt.Sprintf("get oidc provider: %s", string(errBody)))
	}

	var getOidcProviderResp getOidcProviderResponse
	err = json.NewDecoder(response.Body).Decode(&getOidcProviderResp)
	if err != nil {
		return nil, err
	}

	return mapToGetOidcProviderResponse(&getOidcProviderResp), nil
}
