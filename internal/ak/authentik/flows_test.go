package authentik

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/svetlyopet/authentik-cli/internal/ak"
)

func Test_GetFlowBySlug(t *testing.T) {
	const path = "GET /api/v3/flows/instances/default-invalidation-flow/"

	t.Run("found", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			path: respond(http.StatusOK, flowFixture("flow-uuid", "Logout", "default-invalidation-flow", "invalidation")),
		})

		got, err := r.GetFlowBySlug("default-invalidation-flow")
		if err != nil {
			t.Fatalf("GetFlowBySlug() error = %v", err)
		}
		want := &ak.Flow{PK: "flow-uuid", Name: "Logout", Slug: "default-invalidation-flow", Designation: "invalidation"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetFlowBySlug() got = %+v, want %+v", got, want)
		}
	})

	t.Run("not found", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			path: respond(http.StatusNotFound, errorBody("No Flow matches the given query.")),
		})

		_, err := r.GetFlowBySlug("default-invalidation-flow")
		assertNotExists(t, err)
	})
}
