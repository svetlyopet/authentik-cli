package authentik

import (
	"github.com/svetlyopet/authentik-cli/internal/ak"
)

func mapToCreateOrUpdateRoleResponse(role *createOrUpdateRoleResponse) *ak.Role {
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

func mapToCreateOrUpdateGroupResponse(group *createOrUpdateGroupResponse) *ak.Group {
	return &ak.Group{
		PK:   group.PK,
		Name: group.Name,
		GroupAttributes: ak.GroupAttributes{
			Tenant: group.Attributes.Tenant,
		},
	}
}

func mapToGetGroupByNameResponse(roles *getGroupsResponse) *ak.Group {
	res := &ak.Group{}
	for _, role := range roles.Results {
		res = &ak.Group{
			PK:   role.PK,
			Name: role.Name,
			GroupAttributes: ak.GroupAttributes{
				Tenant: role.Attributes.Tenant,
			},
		}
	}

	return res
}

func mapToCreateOrUpdateUserResponse(user *createOrUpdateUserResponse) *ak.User {
	return &ak.User{
		PK:       user.PK,
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Path:     user.Path,
		IsActive: user.IsActive,
		Attributes: ak.UserAttributes{
			UserType: user.Attributes.UserType,
			Tenant:   user.Attributes.Tenant,
		},
	}
}

func mapToUserGetResponse(user *getUserResponse) *ak.User {
	var userGetResponse ak.User

	for _, userResults := range user.Results {
		userGetResponse.PK = userResults.PK
		userGetResponse.Username = userResults.Username
		userGetResponse.Name = userResults.Name
		userGetResponse.Email = userResults.Email
		userGetResponse.Path = userResults.Path
		userGetResponse.IsActive = userResults.IsActive
		userGetResponse.Attributes.UserType = userResults.Attributes.UserType
	}

	return &userGetResponse
}
