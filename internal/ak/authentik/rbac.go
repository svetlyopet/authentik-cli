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
	rbacRolePath                  = "%s/api/v3/rbac/roles/"
	rbacRolePathUpdateDelete      = "%s/api/v3/rbac/roles/%s/"
	rbacRolePermissionsAssignPath = "%s/api/v3/rbac/permissions/assigned_by_roles/%s/assign/"
)

func (a *authentik) CreateRole(name string) (*ak.Role, error) {
	createRoleRequest := createRoleRequest{
		Name: name,
	}

	createRoleRequestBytes, err := json.Marshal(createRoleRequest)
	if err != nil {
		return nil, err
	}

	response, err := a.doRequest(http.MethodPost,
		fmt.Sprintf(rbacRolePath, a.url),
		bytes.NewBuffer(createRoleRequestBytes))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(response.Body)
		return nil, customErrors.NewUnexpectedResult(fmt.Sprintf("create role: %s", string(errBody)))
	}

	var createRoleResp createRoleResponse
	err = json.NewDecoder(response.Body).Decode(&createRoleResp)
	if err != nil {
		return nil, err
	}

	return mapToCreateOrUpdateRoleResponse(&createRoleResp), nil
}

func (a *authentik) GetRoleByName(name string) (*ak.Role, error) {
	param := url.Values{}
	param.Add("search", name)

	response, err := a.doRequestWithQuery(http.MethodGet, fmt.Sprintf(rbacRolePath, a.url), nil, &param)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("get role by name: %s", string(errBody))
	}

	var getRolesResp getRolesResponse
	err = json.NewDecoder(response.Body).Decode(&getRolesResp)
	if err != nil {
		return nil, err
	}

	if len(getRolesResp.Results) == 0 {
		return nil, customErrors.NewNotExists("get role by name: role not found")
	}

	if len(getRolesResp.Results) > 1 {
		return nil, customErrors.NewUnexpectedResult("get role by name: found more than one role with the search query")
	}

	return mapToGetRoleByNameResponse(&getRolesResp), nil
}

func (a *authentik) AssignTenantAdminPermissionsToRole(rolePK string) error {
	assignPermissionsRequest := assignPermissionsRequest{
		Permissions: []string{
			AccessAdminInterfacePerm,
			ViewSystemInfoPerm,
			ViewBrandPerm,
			ViewOutpostPerm,
			ViewApplicationPerm,
			ViewProviderPerm,
			ViewEventPerm,
		},
	}

	assignPermissionsRequestBytes, err := json.Marshal(assignPermissionsRequest)
	if err != nil {
		return err
	}

	response, err := a.doRequest(http.MethodPost,
		fmt.Sprintf(rbacRolePermissionsAssignPath, a.url, rolePK),
		bytes.NewBuffer(assignPermissionsRequestBytes))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(response.Body)
		return customErrors.NewUnexpectedResult(fmt.Sprintf("assign permissions to role: %s", string(errBody)))
	}

	return nil
}

func (a *authentik) DeleteRole(uuid string) error {
	response, err := a.doRequest(http.MethodDelete, fmt.Sprintf(rbacRolePathUpdateDelete, a.url, uuid), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		errBody, _ := io.ReadAll(response.Body)
		return customErrors.NewUnexpectedResult(fmt.Sprintf("delete role: %s", string(errBody)))
	}

	return nil
}
