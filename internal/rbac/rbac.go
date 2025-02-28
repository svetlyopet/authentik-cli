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

	fmt.Printf("created role %s\n", name)

	return role, nil
}
