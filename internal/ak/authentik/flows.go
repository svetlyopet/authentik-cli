package authentik

import (
	"context"

	"github.com/svetlyopet/authentik-cli/internal/ak"
)

// GetFlowBySlug returns the flow with the given slug, or a NotExists error
// when authentik has no such flow.
func (a *authentik) GetFlowBySlug(slug string) (*ak.Flow, error) {
	if err := a.preflightCheck(); err != nil {
		return nil, err
	}

	f, resp, err := a.client.FlowsAPI.FlowsInstancesRetrieve(context.Background(), slug).Execute()
	if err != nil {
		return nil, apiError("get flow "+slug, resp, err)
	}

	flow := mapFlow(f)
	return &flow, nil
}
