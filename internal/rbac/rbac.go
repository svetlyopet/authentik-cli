package rbac

import (
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

func CreateRole(tenantName string) (rolePK string, err error) {
	roleName := fmt.Sprintf(constants.TenantAdminRbacRoleNamePattern, tenantName)

	role, err := ak.Repo.CreateRole(roleName)
	if err != nil {
		return "", err
	}

	fmt.Printf("created role %s\n", roleName)

	return role.PK, nil
}
