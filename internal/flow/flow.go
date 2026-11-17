package flow

import "github.com/svetlyopet/authentik-cli/internal/ak"

func GetFlowBySlug(slug string) (*ak.Flow, error) {
	return ak.Repo.GetFlowBySlug(slug)
}
