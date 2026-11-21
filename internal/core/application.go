package core

import (
	"errors"
	"strings"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/diff"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	"github.com/svetlyopet/authentik-cli/internal/logger"
)

func CreateApplication(name, slug string, providerPK int) (err error) {
	if slug == "" {
		slug = strings.ToLower(name)
	}

	_, err = ak.Repo.CreateApplication(name, slug, providerPK)
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeApplication, constants.ActionCreated, name)

	return nil
}

func GetApplication(name string) (application *ak.Application, err error) {
	application, err = ak.Repo.GetApplicationByName(name)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func DeleteApplication(name, slug string) (err error) {
	err = ak.Repo.DeleteApplication(slug)
	if err != nil {
		var notExistsError *customErrors.NotExists
		if errors.As(err, &notExistsError) {
			return nil
		}
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeApplication, constants.ActionDeleted, name)

	return nil
}

// ApplicationChanges lists the changes requested for an application. A nil
// field is left unchanged.
type ApplicationChanges struct {
	Name *string
	Slug *string
}

// UpdateApplication applies the requested changes to application. Fields that
// were not requested, or already hold the requested value, are not sent to
// authentik and keep their current value.
func UpdateApplication(application *ak.Application, changes ApplicationChanges) error {
	patch := ak.ApplicationPatch{
		Name: diff.Changed(application.Name, changes.Name),
		Slug: diff.Changed(application.Slug, changes.Slug),
	}
	if patch.IsEmpty() {
		logger.LogObjectChange(constants.ObjectTypeApplication, constants.ActionUnchanged, application.Name)
		return nil
	}

	updated, err := ak.Repo.UpdateApplication(application.Slug, patch)
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeApplication, constants.ActionChanged, updated.Name)

	return nil
}
