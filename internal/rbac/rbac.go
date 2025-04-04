package rbac

import (
	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
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

func DeleteRole(name, uuid string) (err error) {
	err = ak.Repo.DeleteRole(uuid)
	if err != nil {
		return err
	}

	logger.WriteStdout(constants.ObjectTypeRole, constants.ActionDeleted, name)

	return nil
}
