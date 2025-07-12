package core

import "github.com/svetlyopet/authentik-cli/internal/ak"

func mapToGetGroupDetails(group *ak.Group) *Group {
	var roles = []string{}
	for _, role := range group.Roles {
		roles = append(roles, role.Name)
	}

	return &Group{
		Name:   group.Name,
		Tenant: group.GroupAttributes.Tenant,
		Roles:  roles,
	}
}

func mapToGetUserDetails(user *ak.User) *User {
	return &User{
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
	}
}
