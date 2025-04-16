package tenant

import (
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
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
		logger.LogObjectChange(constants.ObjectTypeRole, constants.ActionUnchanged, roleName)
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
		logger.LogObjectChange(constants.ObjectTypeGroup, constants.ActionUnchanged, groupName)
	}

	return nil
}

func Delete(name string) (err error) {
	roleName := fmt.Sprintf(constants.TenantAdminRbacRoleNamePattern, name)

	err = rbac.DeleteRole(roleName)
	if err != nil {
		return err
	}

	groupName := fmt.Sprintf(constants.TenantAdminGroupNamePattern, name)

	err = core.DeleteGroup(groupName)
	if err != nil {
		return err
	}

	return nil
}
