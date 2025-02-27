package tenant

import "github.com/svetlyopet/authentik-cli/internal/rbac"

func Create(name string) error {
	_, err := rbac.CreateRole(name)
	if err != nil {
		return err
	}

	return nil
}

func Delete(name string) error {
	return nil
}
