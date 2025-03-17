package core

import (
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
)

func CreateGroup(name string, roles []string, attributes map[string]string) (*ak.Group, error) {
	group, err := ak.Repo.CreateGroup(name, roles, attributes)
	if err != nil {
		return nil, err
	}

	fmt.Printf("group/%s created\n", name)

	return group, nil
}

func GetGroupByName(name string) (group *ak.Group, err error) {
	group, err = ak.Repo.GetGroupByName(name)
	if err != nil {
		return nil, err
	}

	return group, nil
}
