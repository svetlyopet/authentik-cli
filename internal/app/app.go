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

// Changes lists the changes requested for an application and its provider.
// A nil field is left unchanged.
type Changes struct {
	Name *string
	Slug *string
	Oidc provider.OidcChanges
}

// Update applies the requested changes to the application with the given
// name and to its provider. Anything that was not requested keeps its current
// value.
func Update(name string, changes Changes) error {
	application, err := core.GetApplication(name)
	if err != nil {
		return err
	}

	if !changes.Oidc.IsEmpty() && application.ProviderType != constants.ProviderTypeOIDC {
		return fmt.Errorf("application %s has no OIDC provider, so OIDC settings cannot be changed", name)
	}

	if application.ProviderType == constants.ProviderTypeOIDC {
		oidcChanges := changes.Oidc
		// Providers are created with their application's name; keep the two
		// in step when the application is renamed.
		if changes.Name != nil && application.ProviderName == application.Name {
			oidcChanges.Name = changes.Name
		}

		if !oidcChanges.IsEmpty() {
			if err := provider.UpdateOidcProvider(application.ProviderPK, oidcChanges); err != nil {
				return err
			}
		}
	}

	if changes.Name == nil && changes.Slug == nil {
		return nil
	}

	return core.UpdateApplication(application, core.ApplicationChanges{
		Name: changes.Name,
		Slug: changes.Slug,
	})
}
