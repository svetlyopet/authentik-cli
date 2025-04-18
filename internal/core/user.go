package core

import (
	"errors"
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	"github.com/svetlyopet/authentik-cli/internal/logger"
)

func CreateUser(username, name, surname, email, tenant string) (err error) {
	userPath := "users"
	userAttr := ak.UserAttributes{}

	if surname != "" {
		name = fmt.Sprintf("%s %s", name, surname)
	}

	tenantGroup := &ak.Group{}
	if tenant != "" {
		tenantGroup, err = GetGroupByName(fmt.Sprintf(constants.TenantAdminGroupNamePattern, tenant))
		if err != nil {
			return customErrors.NewNotExists("tenant not found")
		}
		userPath = "tenant-admins"
		userAttr.UserType = "tenant-admin"
		userAttr.Tenant = tenant
	}

	usr := ak.User{
		Username:   username,
		Name:       name,
		Email:      email,
		Path:       userPath,
		IsActive:   true,
		Attributes: userAttr,
	}

	user, err := ak.Repo.CreateUser(usr)
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeUser, constants.ActionCreated, username)

	if tenant != "" {
		err = AddUserToGroup(user.PK, tenantGroup.PK, tenantGroup.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteUser(username string) (err error) {
	user, err := ak.Repo.GetUserByUsername(username)
	if err != nil {
		var notExistsError *customErrors.NotExists
		if errors.As(err, &notExistsError) {
			return nil
		}
	}

	err = ak.Repo.DeleteUser(fmt.Sprintf("%d", user.PK))
	if err != nil {
		return err
	}

	logger.LogObjectChange(constants.ObjectTypeUser, constants.ActionDeleted, username)

	return nil
}
