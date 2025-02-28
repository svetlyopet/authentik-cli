package tenant

import (
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
	"github.com/svetlyopet/authentik-cli/internal/rbac"
)

func Create(name string) error {
	attributes := map[string]string{}
	attributes["tenant"] = name

	roleName := fmt.Sprintf(constants.TenantAdminRbacRoleNamePattern, name)
	role, err := rbac.CreateRole(roleName)
	if err != nil {
		return err
	}

	groupName := fmt.Sprintf(constants.TenantAdminGroupNamePattern, name)
	_, err = core.CreateGroup(groupName, []string{role.PK}, attributes)
	if err != nil {
		return err
	}

	return nil
}

func Delete(name string) error {
	return nil
}
