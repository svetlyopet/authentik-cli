package core

import (
	"github.com/svetlyopet/authentik-cli/internal/ak"
)

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
	var groups = []string{}
	for _, group := range user.Groups {
		groups = append(groups, group.Name)
	}

	return &User{
		Name:        user.Name,
		Email:       user.Email,
		IsActive:    user.IsActive,
		IsSuperuser: user.IsSuperuser,
		Attributes:  user.Attributes,
		Groups:      groups,
	}
}
