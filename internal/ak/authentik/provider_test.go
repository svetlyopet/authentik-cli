package authentik

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/svetlyopet/authentik-cli/internal/ak"
)

func Test_CreateOidcProvider(t *testing.T) {
	r := newFakeAuthentik(t, map[string]route{
		"POST /api/v3/providers/oauth2/": func(t *testing.T, r *http.Request) (int, any) {
			body := decodeBody(t, r)
			for key, want := range map[string]any{
				"name":               "example-provider",
				"authorization_flow": "authz-flow",
				"invalidation_flow":  "invalidation-flow",
				"client_type":        "confidential",
				"sub_mode":           "hashed_user_id",
				"issuer_mode":        "per_provider",
			} {
				if body[key] != want {
					t.Errorf("request %s = %v, want %v", key, body[key], want)
				}
			}
			if _, ok := body["authentication_flow"]; ok {
				t.Error("authentication_flow should be omitted when empty")
			}
			uris, _ := body["redirect_uris"].([]any)
			if len(uris) != 1 {
				t.Fatalf("redirect_uris = %v, want one entry", body["redirect_uris"])
			}

			return http.StatusCreated, oauth2ProviderFixture(9, "example-provider", obj{
				"authorization_flow": "authz-flow",
				"invalidation_flow":  "invalidation-flow",
				"client_type":        "confidential",
				"client_id":          "client-id",
				"redirect_uris":      list(obj{"matching_mode": "strict", "url": "https://app.example.com/cb"}),
				"sub_mode":           "hashed_user_id",
				"issuer_mode":        "per_provider",
			})
		},
	})

	got, err := r.CreateOidcProvider(ak.OidcProvider{
		Provider: ak.Provider{
			Name:              "example-provider",
			AuthorizationFlow: "authz-flow",
			InvalidationFlow:  "invalidation-flow",
		},
		ClientType:   "confidential",
		RedirectUris: []ak.OidcRedirectUri{{MatchingMode: "strict", Url: "https://app.example.com/cb"}},
		SubMode:      "hashed_user_id",
		IssuerMode:   "per_provider",
	})
	if err != nil {
		t.Fatalf("CreateOidcProvider() error = %v", err)
	}

	want := &ak.OidcProvider{
		Provider: ak.Provider{
			PK:                9,
			Name:              "example-provider",
			AuthorizationFlow: "authz-flow",
			InvalidationFlow:  "invalidation-flow",
		},
		ClientType:   "confidential",
		ClientId:     "client-id",
		RedirectUris: []ak.OidcRedirectUri{{MatchingMode: "strict", Url: "https://app.example.com/cb"}},
		SubMode:      "hashed_user_id",
		IssuerMode:   "per_provider",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateOidcProvider() got = %+v, want %+v", got, want)
	}
}

func Test_GetOidcProvider(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"GET /api/v3/providers/oauth2/9/": respond(http.StatusOK, oauth2ProviderFixture(9, "example-provider", nil)),
		})

		got, err := r.GetOidcProvider(9)
		if err != nil {
			t.Fatalf("GetOidcProvider() error = %v", err)
		}
		if got.PK != 9 || got.Name != "example-provider" {
			t.Errorf("GetOidcProvider() got = %+v", got)
		}
	})

	t.Run("not found", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"GET /api/v3/providers/oauth2/9/": respond(http.StatusNotFound, errorBody("not found")),
		})

		_, err := r.GetOidcProvider(9)
		assertNotExists(t, err)
	})
}

func Test_DeleteProvider(t *testing.T) {
	t.Run("deleted", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{"DELETE /api/v3/providers/all/9/": respond(http.StatusNoContent, nil)})
		if err := r.DeleteProvider(9); err != nil {
			t.Errorf("DeleteProvider() error = %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"DELETE /api/v3/providers/all/9/": respond(http.StatusNotFound, errorBody("not found")),
		})
		assertNotExists(t, r.DeleteProvider(9))
	})
}
