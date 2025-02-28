package core

import (
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/pkg/idp"
)

func CreateGroup(name string, roles []string, attributes map[string]string) (*idp.Group, error) {
	group, err := ak.Repo.CreateGroup(name, roles, attributes)
	if err != nil {
		return nil, err
	}

	fmt.Printf("created group %s\n", name)

	return group, nil
}
