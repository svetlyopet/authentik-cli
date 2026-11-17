package authentik

import (
	"context"

	api "goauthentik.io/api/v3"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
)

// tenantRoleViewPermissions are granted to every tenant role so its members
// can browse the admin interface.
var tenantRoleViewPermissions = []string{
	AccessAdminInterfacePerm,
	ViewSystemInfoPerm,
	ViewBrandPerm,
	ViewFlowPerm,
	ViewOutpostPerm,
	ViewApplicationPerm,
	ViewProviderPerm,
	ViewEventPerm,
}

func (a *authentik) CreateRole(name string) (*ak.Role, error) {
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	role, resp, err := a.client.RbacAPI.RbacRolesCreate(context.Background()).
		RoleRequest(*api.NewRoleRequest(name)).
		Execute()
	if err != nil {
		return nil, apiError("create role", resp, err)
	}

	return mapRole(role), nil
}

func (a *authentik) GetRoleByName(name string) (*ak.Role, error) {
	const op = "get role by name"
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	// The roles endpoint only offers a fuzzy search, so narrow the results
	// down to exact name matches.
	roles, resp, err := a.client.RbacAPI.RbacRolesList(context.Background()).
		Search(name).
		Execute()
	if err != nil {
		return nil, apiError(op, resp, err)
	}

	var matches []api.Role
	for _, role := range roles.Results {
		if role.Name == name {
			matches = append(matches, role)
		}
	}

	switch len(matches) {
	case 0:
		return nil, customErrors.NewNotExists(op + ": role not found")
	case 1:
		return mapRole(&matches[0]), nil
	default:
		return nil, customErrors.NewUnexpectedResult(op + ": found more than one role with the search query")
	}
}

func (a *authentik) AssignViewPermissionsToTenantRole(rolePK string) error {
	if err := a.preflightCheck(); err != nil {
		return err
	}

	_, resp, err := a.client.RbacAPI.RbacPermissionsAssignedByRolesAssign(context.Background(), rolePK).
		PermissionAssignRequest(*api.NewPermissionAssignRequest(tenantRoleViewPermissions)).
		Execute()
	if err != nil {
		return apiError("assign permissions to role", resp, err)
	}

	return nil
}

func (a *authentik) DeleteRole(uuid string) error {
	if err := a.preflightCheck(); err != nil {
		return err
	}

	resp, err := a.client.RbacAPI.RbacRolesDestroy(context.Background(), uuid).Execute()
	if err != nil {
		return apiError("delete role", resp, err)
	}

	return nil
}
