package authentik

import "github.com/svetlyopet/authentik-cli/pkg/idp"

func mapToCreateOrUpdateRoleResponse(role *createRoleResponse) *idp.Role {
	return &idp.Role{
		PK:   role.PK,
		Name: role.Name,
	}
}

func mapToCreateOrUpdateGroupResponse(group *createGroupResponse) *idp.Group {
	return &idp.Group{
		PK:   group.PK,
		Name: group.Name,
	}
}
