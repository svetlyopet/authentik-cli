package tenant

import (
	"errors"
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	"github.com/svetlyopet/authentik-cli/internal/logger"
	"github.com/svetlyopet/authentik-cli/internal/rbac"
)

func Create(name string) (err error) {
	role := &ak.Role{}

	roleName := fmt.Sprintf(constants.TenantAdminRbacRoleNamePattern, name)

	if role, err = rbac.GetRoleByName(roleName); err != nil {
		role, err = rbac.CreateRole(roleName)
		if err != nil {
			return err
		}
	} else {
		logger.WriteStdout(constants.ObjectTypeRole, constants.ActionUnchanged, roleName)
	}

	groupName := fmt.Sprintf(constants.TenantAdminGroupNamePattern, name)
	groupAttributes := ak.GroupAttributes{
		Tenant: name,
	}

	if _, err = core.GetGroupByName(groupName); err != nil {
		_, err = core.CreateGroup(groupName, []string{role.PK}, groupAttributes)
		if err != nil {
			return err
		}
	} else {
		logger.WriteStdout(constants.ObjectTypeGroup, constants.ActionUnchanged, groupName)
	}

	return nil
}

func Delete(name string) (err error) {
	roleName := fmt.Sprintf(constants.TenantAdminRbacRoleNamePattern, name)

	var notExistsError *customErrors.NotExists
	role, err := rbac.GetRoleByName(roleName)
	if err != nil {
		if !errors.As(err, &notExistsError) {
			return err
		}
	}

	if role != nil {
		err = rbac.DeleteRole(role.Name, role.PK)
		if err != nil {
			return err
		}
	}

	groupName := fmt.Sprintf(constants.TenantAdminGroupNamePattern, name)

	group, err := core.GetGroupByName(groupName)
	if err != nil {
		if !errors.As(err, &notExistsError) {
			return err
		}
	}
	if group != nil {
		err = core.DeleteGroup(group.Name, group.PK)
		if err != nil {
			return err
		}
	}

	return nil
}
