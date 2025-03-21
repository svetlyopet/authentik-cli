package authentik

import (
	"github.com/svetlyopet/authentik-cli/internal/ak"
)

func mapToCreateOrUpdateRoleResponse(role *createRoleResponse) *ak.Role {
	return &ak.Role{
		PK:   role.PK,
		Name: role.Name,
	}
}

func mapToGetRoleByNameResponse(roles *getRolesResponse) *ak.Role {
	res := &ak.Role{}
	for _, role := range roles.Results {
		res = &ak.Role{
			PK:   role.PK,
			Name: role.Name,
		}
	}

	return res
}

func mapToCreateOrUpdateGroupResponse(group *createGroupResponse) *ak.Group {
	return &ak.Group{
		PK:   group.PK,
		Name: group.Name,
	}
}

func mapToGetGroupByNameResponse(roles *getGroupsResponse) *ak.Group {
	res := &ak.Group{}
	for _, role := range roles.Results {
		res = &ak.Group{
			PK:   role.PK,
			Name: role.Name,
		}
	}

	return res
}
