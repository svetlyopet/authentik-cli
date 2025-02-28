package authentik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/svetlyopet/authentik-cli/pkg/idp"
)

const coreGroupPath = "%s/api/v3/core/groups/"

func (a *authentik) CreateGroup(name string, roles []string, attributes map[string]string) (*idp.Group, error) {
	createGroupRequest := createGroupRequest{
		Name:        name,
		IsSuperuser: false,
		Users:       []int{},
		Attributes:  attributes,
		Roles:       roles,
	}

	createGroupRequestBytes, err := json.Marshal(createGroupRequest)
	if err != nil {
		return nil, err
	}

	response, err := a.doRequest(http.MethodPost,
		fmt.Sprintf(coreGroupPath, a.url),
		bytes.NewBuffer(createGroupRequestBytes))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("create group: %s", string(errBody))
	}

	var createGroupResp createGroupResponse
	err = json.NewDecoder(response.Body).Decode(&createGroupResp)
	if err != nil {
		return nil, err
	}

	return mapToCreateOrUpdateGroupResponse(&createGroupResp), nil
}

func (a *authentik) GetGroupByName(name string) (*idp.Group, error) {
	param := url.Values{}
	param.Add("name", name)

	response, err := a.doRequestWithQuery(http.MethodGet, fmt.Sprintf(coreGroupPath, a.url), nil, &param)
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
		return nil, fmt.Errorf("get group by name: group not found")
	}

	if len(getGroupsResp.Results) > 1 {
		return nil, fmt.Errorf("get group by name: found more than one group with the search query")
	}

	return mapToGetGroupByNameResponse(&getGroupsResp), nil
}
