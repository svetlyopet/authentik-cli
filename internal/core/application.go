package core

import (
	"errors"
	"strings"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
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

	logger.WriteStdout(constants.ObjectTypeApplication, constants.ActionCreated, name)

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

	logger.WriteStdout(constants.ObjectTypeApplication, constants.ActionDeleted, name)

	return nil
}
