package authentik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/svetlyopet/authentik-cli/pkg/idp"
)

const rbacRolePath = "%s/api/v3/rbac/roles/"

func (a *authentik) CreateRole(name string) (*idp.Role, error) {
	createRoleRequest := createRoleRequest{
		Name: name,
	}

	createRoleRequestBytes, err := json.Marshal(createRoleRequest)
	if err != nil {
		return nil, err
	}

	response, err := a.doRequest(http.MethodPost, fmt.Sprintf(rbacRolePath, a.url), bytes.NewBuffer(createRoleRequestBytes))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("create role: %s", string(errBody))
	}

	var createRoleResp createRoleResponse
	err = json.NewDecoder(response.Body).Decode(&createRoleResp)
	if err != nil {
		return nil, err
	}

	return mapToCreateOrUpdateRoleResponse(&createRoleResp), nil
}
