package app

import (
	"errors"
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	"github.com/svetlyopet/authentik-cli/internal/provider"
)

func CreateOidc(name, slug, oidcClientType, oidcConsentType string, oidcEncryptToken bool, oidcRedirectUris []string) error {
	var providerPK int

	oidcProvider, err := provider.CreateOidcProvider(name, oidcClientType, oidcConsentType, oidcEncryptToken, oidcRedirectUris)
	if err != nil {
		return err
	}
	providerPK = oidcProvider.PK

	err = core.CreateApplication(name, slug, providerPK)
	if err != nil {
		return err
	}

	return nil
}

func Get(name string) (*App, error) {
	var appDetails = &App{}

	application, err := core.GetApplication(name)
	if err != nil {
		return nil, err
	}

	switch application.ProviderType {
	case constants.ProviderTypeOIDC:
		prov, err := provider.GetOidcProvider(application.ProviderPK)
		if err != nil {
			return nil, err
		}

		appDetails = mapToGetAppWithOidcProvider(application, *prov)
	case constants.ProviderTypeLDAP:
		return nil, fmt.Errorf("LDAP provider is not supported yet")
	case constants.ProviderTypeSAML:
		return nil, fmt.Errorf("SAML provider is not supported yet")
	case constants.ProviderTypeProxy:
		return nil, fmt.Errorf("Proxy provider is not supported yet")
	case constants.ProviderTypeRAC:
		return nil, fmt.Errorf("RAC provider is not supported yet")
	case constants.ProviderTypeSCIM:
		return nil, fmt.Errorf("SCIM provider is not supported yet")
	default:
		return nil, fmt.Errorf("%s is not a supported provider type", application.ProviderType)
	}

	return appDetails, nil
}

func Delete(name string) error {
	app, err := core.GetApplication(name)
	if err != nil {
		var notExistsError *customErrors.NotExists
		if errors.As(err, &notExistsError) {
			return nil
		}
	}

	if err = core.DeleteApplication(app.Name, app.Slug); err != nil {
		var notExistsError *customErrors.NotExists
		if errors.As(err, &notExistsError) {
			return nil
		}
	}

	if err = provider.DeleteProvider(app.Name, app.ProviderPK); err != nil {
		return err
	}

	return nil
}
