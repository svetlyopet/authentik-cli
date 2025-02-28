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

	fmt.Printf("group/%s created\n", name)

	return group, nil
}

func GetGroupByName(name string) (group *idp.Group, err error) {
	group, err = ak.Repo.GetGroupByName(name)
	if err != nil {
		return nil, err
	}

	return group, nil
}
