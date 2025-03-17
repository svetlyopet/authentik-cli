package tenant

import (
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
	"github.com/svetlyopet/authentik-cli/internal/rbac"
)

func Create(name string) (err error) {
	role := &ak.Role{}

	attributes := map[string]string{}
	attributes["tenant"] = name

	roleName := fmt.Sprintf(constants.TenantAdminRbacRoleNamePattern, name)

	if role, err = rbac.GetRoleByName(roleName); err != nil {
		role, err = rbac.CreateRole(roleName)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("role/%s unchanged\n", roleName)
	}

	groupName := fmt.Sprintf(constants.TenantAdminGroupNamePattern, name)

	if _, err = core.GetGroupByName(groupName); err != nil {
		_, err = core.CreateGroup(groupName, []string{role.PK}, attributes)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("group/%s unchanged\n", groupName)
	}

	return nil
}

func Delete(name string) error {
	return nil
}
