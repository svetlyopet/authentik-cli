package core

import (
	"strings"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
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
