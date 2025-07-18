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
	coreUsersPath                    = "%s/api/v3/core/users/"
	coreUsersPathUpdateDeletePath    = "%s/api/v3/core/users/%s/"
	coreGroupsPath                   = "%s/api/v3/core/groups/"
	coreGroupsPathUpdateDeletePath   = "%s/api/v3/core/groups/%s/"
	coreGroupsAddUserPath            = "%s/api/v3/core/groups/%s/add_user/"
	coreApplicationsPath             = "%s/api/v3/core/applications/"
	coreApplicationsUpdateDeletePath = "%s/api/v3/core/applications/%s/"
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
	defer response.Body.Close() //nolint

	if response.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(response.Body)
		return nil, customErrors.NewUnexpectedResult(fmt.Sprintf("create group: %s", string(errBody)))
	}

	var createGroupResp getGroupResponse
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
	defer response.Body.Close() //nolint

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

func (a *authentik) GetGroup(uuid string) (*ak.Group, error) {
	response, err := a.doRequest(http.MethodGet, fmt.Sprintf(coreGroupsPathUpdateDeletePath, a.url, uuid), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() //nolint

	if response.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("get group: %s", string(errBody))
	}

	var getGroupResp getGroupResponse
	err = json.NewDecoder(response.Body).Decode(&getGroupResp)
	if err != nil {
		return nil, err
	}

	return mapToGetGroupResponse(&getGroupResp), nil
}

func (a *authentik) DeleteGroup(uuid string) error {
	response, err := a.doRequest(http.MethodDelete, fmt.Sprintf(coreGroupsPathUpdateDeletePath, a.url, uuid), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close() //nolint

	if response.StatusCode != http.StatusNoContent {
		errBody, _ := io.ReadAll(response.Body)
		returnErrBody := fmt.Sprintf("delete group: %s", string(errBody))
		return customErrors.NewUnexpectedResult(returnErrBody)
	}

	return nil
}

func (a *authentik) CreateUser(usr ak.User) (*ak.User, error) {
	createUserReq := createUserRequest{
		Username:    usr.Username,
		Name:        usr.Name,
		Email:       usr.Email,
		Path:        usr.Path,
		IsActive:    usr.IsActive,
		IsSuperuser: usr.IsSuperuser,
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

	defer response.Body.Close() //nolint
	if response.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(response.Body)
		returnErrBody := fmt.Sprintf("create user: %s", errBody)
		return nil, customErrors.NewUnexpectedResult(returnErrBody)
	}

	var userResp getUserResponse
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

	defer response.Body.Close() //nolint
	if response.StatusCode != http.StatusNoContent {
		errBody, _ := io.ReadAll(response.Body)
		returnErrBody := fmt.Sprintf("add user to group: %s", errBody)
		return customErrors.NewUnexpectedResult(returnErrBody)
	}

	return nil
}

func (a *authentik) GetUserByUsername(username string) (*ak.User, error) {
	params := url.Values{}
	params.Add("username", username)

	response, err := a.doRequestWithQuery(http.MethodGet, fmt.Sprintf(coreUsersPath, a.url), nil, &params)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close() //nolint

	if response.StatusCode == http.StatusNotFound {
		return nil, customErrors.NewNotExists("user not found")
	}
	if response.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(response.Body)
		returnErrBody := fmt.Sprintf("get user by username: %s", errBody)
		return nil, customErrors.NewUnexpectedResult(returnErrBody)
	}

	var usersResp getUsersResponse
	err = json.NewDecoder(response.Body).Decode(&usersResp)
	if err != nil {
		return nil, err
	}

	if len(usersResp.Results) == 0 {
		return nil, customErrors.NewNotExists("user not found")
	}

	if len(usersResp.Results) > 1 {
		return nil, customErrors.NewUnexpectedResult(fmt.Sprintf("found more than 1 user while searching for %s", username))
	}
	return mapToUsersGetResponse(&usersResp), nil
}

func (a *authentik) DeleteUser(userPK string) error {
	response, err := a.doRequest(http.MethodDelete, fmt.Sprintf(coreUsersPathUpdateDeletePath, a.url, userPK), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close() //nolint

	if response.StatusCode != http.StatusNoContent {
		errBody, _ := io.ReadAll(response.Body)
		returnErrBody := fmt.Sprintf("delete user: %s", string(errBody))
		return customErrors.NewUnexpectedResult(returnErrBody)
	}

	return nil
}

func (a *authentik) CreateApplication(name, slug string, providerPK int) (*ak.Application, error) {
	createApplicationRequest := createOrUpdateApplicationRequest{
		Name:     name,
		Slug:     slug,
		Provider: providerPK,
	}

	createApplicationRequestBytes, err := json.Marshal(createApplicationRequest)
	if err != nil {
		return nil, err
	}

	response, err := a.doRequest(http.MethodPost,
		fmt.Sprintf(coreApplicationsPath, a.url),
		bytes.NewBuffer(createApplicationRequestBytes))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() //nolint

	if response.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(response.Body)
		returnErrBody := fmt.Sprintf("create application: %s", errBody)
		return nil, customErrors.NewUnexpectedResult(returnErrBody)
	}

	var createApplicationResp getApplicationResponse
	err = json.NewDecoder(response.Body).Decode(&createApplicationResp)
	if err != nil {
		return nil, err
	}

	return mapToCreateOrUpdateApplicationResponse(&createApplicationResp), nil
}

func (a *authentik) GetApplicationByName(name string) (*ak.Application, error) {
	params := url.Values{}
	params.Add("search", name)

	response, err := a.doRequestWithQuery(http.MethodGet, fmt.Sprintf(coreApplicationsPath, a.url), nil, &params)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() //nolint

	if response.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(response.Body)
		returnErrBody := fmt.Sprintf("get application: %s", errBody)
		return nil, customErrors.NewUnexpectedResult(returnErrBody)
	}

	var getApplicationsResp getApplicationsResponse
	err = json.NewDecoder(response.Body).Decode(&getApplicationsResp)
	if err != nil {
		return nil, err
	}

	application, err := mapToGetApplicationsByNameResponse(&getApplicationsResp, name)
	if err != nil {
		return nil, customErrors.NewNotExists("application not found")
	}

	return application, nil
}

func (a *authentik) DeleteApplication(slug string) error {
	response, err := a.doRequest(http.MethodDelete, fmt.Sprintf(coreApplicationsUpdateDeletePath, a.url, slug), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close() //nolint

	if response.StatusCode == http.StatusNotFound {
		return customErrors.NewNotExists("application not found")
	}

	if response.StatusCode != http.StatusNoContent {
		errBody, _ := io.ReadAll(response.Body)
		returnErrBody := fmt.Sprintf("delete application: %s", errBody)
		return customErrors.NewUnexpectedResult(returnErrBody)
	}

	return nil
}
