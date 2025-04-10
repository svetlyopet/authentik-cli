package flow

import "github.com/svetlyopet/authentik-cli/internal/ak"

func GetFlows() (flows []ak.Flow, err error) {
	flows, err = ak.Repo.GetFlows()
	if err != nil {
		return flows, nil
	}

	return flows, nil
}
