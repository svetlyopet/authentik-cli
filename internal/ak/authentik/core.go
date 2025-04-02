package authentik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
)

const (
	coreUsersPath              = "%s/api/v3/core/users/"
	coreGroupsPath             = "%s/api/v3/core/groups/"
	coreGroupsPathUpdateDelete = "%s/api/v3/core/groups/%s/"
	coreGroupsAddUserPath      = "%s/api/v3/core/groups/%s/add_user/"
)

func (a *authentik) CreateGroup(name string, roles []string, attributes ak.GroupAttributes) (*ak.Group, error) {
	createGroupRequest := createGroupRequest{
		Name:        name,
		IsSuperuser: false,
		Users:       []int{},
		Attributes: groupAttributes{
			Tenant: attributes.Tenant,
		},
		Roles: roles,
	}

	createGroupRequestBytes, err := json.Marshal(createGroupRequest)
	if err != nil {
		return nil, err
	}

	response, err := a.doRequest(http.MethodPost,
		fmt.Sprintf(coreGroupsPath, a.url),
		bytes.NewBuffer(createGroupRequestBytes))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(response.Body)
		return nil, customErrors.NewUnexpectedResult(fmt.Sprintf("create group: %s", string(errBody)))
	}

	var createGroupResp createGroupResponse
	err = json.NewDecoder(response.Body).Decode(&createGroupResp)
	if err != nil {
		return nil, err
	}

	return mapToCreateOrUpdateGroupResponse(&createGroupResp), nil
}

func (a *authentik) GetGroupByName(name string) (*ak.Group, error) {
	param := url.Values{}
	param.Add("name", name)

	response, err := a.doRequestWithQuery(http.MethodGet, fmt.Sprintf(coreGroupsPath, a.url), nil, &param)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("get group by name: %s", string(errBody))
	}

	var getGroupsResp getGroupsResponse
	err = json.NewDecoder(response.Body).Decode(&getGroupsResp)
	if err != nil {
		return nil, err
	}

	if len(getGroupsResp.Results) == 0 {
		return nil, customErrors.NewNotExists("get group by name: group not found")
	}

	if len(getGroupsResp.Results) > 1 {
		return nil, customErrors.NewUnexpectedResult("get group by name: found more than one group with the search query")
	}

	return mapToGetGroupByNameResponse(&getGroupsResp), nil
}

func (a *authentik) DeleteGroup(uuid string) error {
	response, err := a.doRequest(http.MethodDelete, fmt.Sprintf(coreGroupsPathUpdateDelete, a.url, uuid), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		errBody, _ := io.ReadAll(response.Body)
		return customErrors.NewUnexpectedResult(fmt.Sprintf("delete group: %s", string(errBody)))
	}

	return nil
}

func (a *authentik) CreateUser(usr ak.User) (*ak.User, error) {
	createUserReq := createUserRequest{
		Username: usr.Username,
		Name:     usr.Name,
		Email:    usr.Email,
		Path:     usr.Path,
		IsActive: usr.IsActive,
		Attributes: userAttributes{
			UserType: usr.Attributes.UserType,
			Tenant:   usr.Attributes.Tenant,
		},
	}

	createUserReqBytes, err := json.Marshal(createUserReq)
	if err != nil {
		return nil, err
	}

	response, err := a.doRequest(http.MethodPost, fmt.Sprintf(coreUsersPath, a.url), bytes.NewBuffer(createUserReqBytes))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(response.Body)
		return nil, customErrors.NewUnexpectedResult(fmt.Sprintf("%s", errBody))
	}

	var userResp createOrUpdateUserResponse
	err = json.NewDecoder(response.Body).Decode(&userResp)
	if err != nil {
		return nil, err
	}

	return mapToCreateOrUpdateUserResponse(&userResp), nil
}

func (a *authentik) AddUserToGroup(userPK int, groupUuid string) error {
	userAddRequest := groupUserAddRequest{
		PK: fmt.Sprintf("%d", userPK),
	}

	userAddRequestBytes, err := json.Marshal(userAddRequest)
	if err != nil {
		return err
	}

	response, err := a.doRequest(http.MethodPost, fmt.Sprintf(coreGroupsAddUserPath, a.url, groupUuid), bytes.NewBuffer(userAddRequestBytes))
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		errBody, _ := io.ReadAll(response.Body)
		return customErrors.NewUnexpectedResult(fmt.Sprintf("%s", errBody))
	}

	return nil
}
