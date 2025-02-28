package authentik

import "github.com/svetlyopet/authentik-cli/pkg/idp"

func mapToCreateOrUpdateRoleResponse(role *createRoleResponse) *idp.Role {
	return &idp.Role{
		PK:   role.PK,
		Name: role.Name,
	}
}

func mapToGetRoleByNameResponse(roles *getRolesResponse) *idp.Role {
	res := &idp.Role{}
	for _, role := range roles.Results {
		res = &idp.Role{
			PK:   role.PK,
			Name: role.Name,
		}
	}

	return res
}

func mapToCreateOrUpdateGroupResponse(group *createGroupResponse) *idp.Group {
	return &idp.Group{
		PK:   group.PK,
		Name: group.Name,
	}
}

func mapToGetGroupByNameResponse(roles *getGroupsResponse) *idp.Group {
	res := &idp.Group{}
	for _, role := range roles.Results {
		res = &idp.Group{
			PK:   role.PK,
			Name: role.Name,
		}
	}

	return res
}
