package authentik

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
)

const (
	flowsInstancesPath = "%s/api/v3/flows/instances/"
)

func (a *authentik) GetFlows() ([]ak.Flow, error) {
	response, err := a.doRequest(http.MethodGet, fmt.Sprintf(flowsInstancesPath, a.url), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() //nolint

	if response.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("get flows: %s", string(errBody))
	}

	var getFlowsResp getFlowsResponse
	err = json.NewDecoder(response.Body).Decode(&getFlowsResp)
	if err != nil {
		return nil, err
	}

	if len(getFlowsResp.Results) == 0 {
		return nil, customErrors.NewNotExists("get flows: no flows found")
	}

	return mapToGetFlowsResponse(&getFlowsResp), nil
}
