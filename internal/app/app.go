package app

import (
	"errors"
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	"github.com/svetlyopet/authentik-cli/internal/provider"
)

func CreateOidc(name, slug, oidcClientType, oidcConsentType string, oidcEncryptToken bool, oidcRedirectUrisStrict, oidcRedirectUrisRegex []string) error {
	oidcProvider, err := provider.CreateOidcProvider(name, oidcClientType, oidcConsentType, oidcEncryptToken, oidcRedirectUrisStrict, oidcRedirectUrisRegex)
	if err != nil {
		return err
	}

	err = core.CreateApplication(name, slug, oidcProvider.PK)
	if err != nil {
		return err
	}

	return nil
}

func Get(name string) (*App, error) {
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

		return mapToGetAppWithOidcProvider(application, *prov), nil
	}

	return nil, fmt.Errorf("%s is not a supported provider type", application.ProviderType)
}

func Delete(name string) error {
	// TODO: also check if a provider exists with the same name and delete it to allow for idempotence
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
