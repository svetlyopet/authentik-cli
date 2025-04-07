package rbac

import (
	"errors"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	"github.com/svetlyopet/authentik-cli/internal/logger"
)

func CreateRole(name string) (role *ak.Role, err error) {
	role, err = ak.Repo.CreateRole(name)
	if err != nil {
		return nil, err
	}

	if err := ak.Repo.AssignViewPermissionsToTenantRole(role.PK); err != nil {
		return nil, err
	}

	logger.WriteStdout(constants.ObjectTypeRole, constants.ActionCreated, name)

	return role, nil
}

func GetRoleByName(name string) (role *ak.Role, err error) {
	role, err = ak.Repo.GetRoleByName(name)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func DeleteRole(name string) (err error) {
	role, err := GetRoleByName(name)
	if err != nil {
		var notExistsError *customErrors.NotExists
		if !errors.As(err, &notExistsError) {
			return err
		}
	}

	err = ak.Repo.DeleteRole(role.PK)
	if err != nil {
		return err
	}

	logger.WriteStdout(constants.ObjectTypeRole, constants.ActionDeleted, name)

	return nil
}
