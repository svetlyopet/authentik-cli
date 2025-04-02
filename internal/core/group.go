package core

import (
	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/logger"
)

func CreateGroup(name string, roles []string, attributes ak.GroupAttributes) (*ak.Group, error) {
	group, err := ak.Repo.CreateGroup(name, roles, attributes)
	if err != nil {
		return nil, err
	}

	logger.WriteStdio(constants.ObjectTypeGroup, constants.ActionCreated, name)

	return group, nil
}

func GetGroupByName(name string) (group *ak.Group, err error) {
	group, err = ak.Repo.GetGroupByName(name)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func DeleteGroup(name, uuid string) (err error) {
	err = ak.Repo.DeleteGroup(uuid)
	if err != nil {
		return err
	}

	logger.WriteStdio(constants.ObjectTypeGroup, constants.ActionDeleted, name)

	return nil
}

func AddUserToGroup(userPK int, groupPK, groupName string) (err error) {
	err = ak.Repo.AddUserToGroup(userPK, groupPK)
	if err != nil {
		return err
	}

	logger.WriteStdio(constants.ObjectTypeGroup, constants.ActionChanged, groupName)

	return nil
}
