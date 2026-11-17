package authentik

import (
	"context"
	"fmt"
	"strconv"

	api "goauthentik.io/api/v3"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
)

func (a *authentik) CreateGroup(name string, roles []string, attributes ak.GroupAttributes) (*ak.Group, error) {
	const op = "create group"
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	req := api.NewGroupRequest(name)
	req.SetIsSuperuser(false)
	req.SetRoles(roles)
	if attrs := groupAttributesToAPI(attributes); len(attrs) > 0 {
		req.SetAttributes(attrs)
	}

	group, resp, err := a.client.CoreAPI.CoreGroupsCreate(context.Background()).GroupRequest(*req).Execute()
	if err != nil {
		return nil, apiError(op, resp, err)
	}

	return mapGroup(group), nil
}

func (a *authentik) GetGroupByName(name string) (*ak.Group, error) {
	const op = "get group by name"
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	groups, resp, err := a.client.CoreAPI.CoreGroupsList(context.Background()).
		Name(name).
		IncludeUsers(false).
		Execute()
	if err != nil {
		return nil, apiError(op, resp, err)
	}

	switch len(groups.Results) {
	case 0:
		return nil, customErrors.NewNotExists(op + ": group not found")
	case 1:
		return mapGroup(&groups.Results[0]), nil
	default:
		return nil, customErrors.NewUnexpectedResult(op + ": found more than one group with the search query")
	}
}

func (a *authentik) GetGroup(uuid string) (*ak.Group, error) {
	const op = "get group"
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	group, resp, err := a.client.CoreAPI.CoreGroupsRetrieve(context.Background(), uuid).
		IncludeUsers(false).
		Execute()
	if err != nil {
		return nil, apiError(op, resp, err)
	}

	return mapGroupWithRoles(group), nil
}

func (a *authentik) DeleteGroup(uuid string) error {
	if err := a.preflightCheck(); err != nil {
		return err
	}

	resp, err := a.client.CoreAPI.CoreGroupsDestroy(context.Background(), uuid).Execute()
	if err != nil {
		return apiError("delete group", resp, err)
	}

	return nil
}

func (a *authentik) CreateUser(usr ak.User) (*ak.User, error) {
	const op = "create user"
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	req := api.NewUserRequest(usr.Username, usr.Name)
	req.SetEmail(usr.Email)
	req.SetPath(usr.Path)
	req.SetIsActive(usr.IsActive)
	if attrs := userAttributesToAPI(usr.Attributes); len(attrs) > 0 {
		req.SetAttributes(attrs)
	}

	user, resp, err := a.client.CoreAPI.CoreUsersCreate(context.Background()).UserRequest(*req).Execute()
	if err != nil {
		return nil, apiError(op, resp, err)
	}

	return mapUser(user), nil
}

func (a *authentik) AddUserToGroup(userPK int, groupUuid string) error {
	const op = "add user to group"
	if err := a.preflightCheck(); err != nil {
		return err
	}

	pk, err := toInt32(op, userPK)
	if err != nil {
		return err
	}

	resp, err := a.client.CoreAPI.CoreGroupsAddUserCreate(context.Background(), groupUuid).
		UserAccountRequest(*api.NewUserAccountRequest(pk)).
		Execute()
	if err != nil {
		return apiError(op, resp, err)
	}

	return nil
}

func (a *authentik) GetUserByUsername(username string) (*ak.User, error) {
	const op = "get user by username"
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	users, resp, err := a.client.CoreAPI.CoreUsersList(context.Background()).
		Username(username).
		Execute()
	if err != nil {
		return nil, apiError(op, resp, err)
	}

	switch len(users.Results) {
	case 0:
		return nil, customErrors.NewNotExists("user not found")
	case 1:
		return mapUser(&users.Results[0]), nil
	default:
		return nil, customErrors.NewUnexpectedResult(fmt.Sprintf("found more than 1 user while searching for %s", username))
	}
}

func (a *authentik) DeleteUser(userPK string) error {
	const op = "delete user"
	if err := a.preflightCheck(); err != nil {
		return err
	}

	id, err := strconv.Atoi(userPK)
	if err != nil {
		return fmt.Errorf("%s: invalid user id %q: %w", op, userPK, err)
	}
	pk, err := toInt32(op, id)
	if err != nil {
		return err
	}

	resp, err := a.client.CoreAPI.CoreUsersDestroy(context.Background(), pk).Execute()
	if err != nil {
		return apiError(op, resp, err)
	}

	return nil
}

func (a *authentik) CreateApplication(name, slug string, providerPK int) (*ak.Application, error) {
	const op = "create application"
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	pk, err := toInt32(op, providerPK)
	if err != nil {
		return nil, err
	}

	req := api.NewApplicationRequest(name, slug)
	req.SetProvider(pk)

	app, resp, err := a.client.CoreAPI.CoreApplicationsCreate(context.Background()).ApplicationRequest(*req).Execute()
	if err != nil {
		return nil, apiError(op, resp, err)
	}

	return mapApplication(app), nil
}

func (a *authentik) GetApplicationByName(name string) (*ak.Application, error) {
	const op = "get application"
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	apps, resp, err := a.client.CoreAPI.CoreApplicationsList(context.Background()).
		Name(name).
		Execute()
	if err != nil {
		return nil, apiError(op, resp, err)
	}

	for i := range apps.Results {
		if apps.Results[i].Name == name {
			return mapApplication(&apps.Results[i]), nil
		}
	}

	return nil, customErrors.NewNotExists("application not found")
}

func (a *authentik) DeleteApplication(slug string) error {
	if err := a.preflightCheck(); err != nil {
		return err
	}

	resp, err := a.client.CoreAPI.CoreApplicationsDestroy(context.Background(), slug).Execute()
	if err != nil {
		return apiError("delete application", resp, err)
	}

	return nil
}
