package rbac

import (
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
)

func CreateRole(name string) (role *ak.Role, err error) {
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

	fmt.Printf("role/%s deleted\n", name)

	return nil
}
