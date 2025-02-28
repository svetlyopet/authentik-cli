package rbac

import (
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/pkg/idp"
)

func CreateRole(name string) (role *idp.Role, err error) {
	role, err = ak.Repo.CreateRole(name)
	if err != nil {
		return nil, err
	}

	if err := ak.Repo.AssignTenantAdminPermissionsToRole(role.PK); err != nil {
		return nil, err
	}

	fmt.Printf("role/%s created\n", name)

	return role, nil
}

func GetRoleByName(name string) (role *idp.Role, err error) {
	role, err = ak.Repo.GetRoleByName(name)
	if err != nil {
		return nil, err
	}

	return role, nil
}
