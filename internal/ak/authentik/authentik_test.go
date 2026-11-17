package authentik

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
)

const testToken = "test-token"

// route answers one request made to the fake authentik server.
type route func(t *testing.T, r *http.Request) (status int, body any)

// respond returns a route replying with status and body (encoded as JSON
// unless nil).
func respond(status int, body any) route {
	return func(*testing.T, *http.Request) (int, any) { return status, body }
}

// errorBody mimics authentik's generic error payload.
func errorBody(detail string) map[string]string {
	return map[string]string{"detail": detail, "code": "error"}
}

// newFakeAuthentik starts an HTTP server that serves routes keyed by
// "METHOD /path" and returns a repository pointed at it. Requests without a
// matching route fail the test.
func newFakeAuthentik(t *testing.T, routes map[string]route) ak.AuthentikRepository {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer "+testToken {
			t.Errorf("Authorization header = %q, want bearer token", got)
		}

		key := r.Method + " " + r.URL.Path
		handle, ok := routes[key]
		if !ok {
			t.Errorf("unexpected request %s", key)
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		status, body := handle(t, r)
		if body != nil {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
	t.Cleanup(srv.Close)

	return New(srv.URL, testToken)
}

// decodeBody decodes the JSON request body into a generic map.
func decodeBody(t *testing.T, r *http.Request) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatalf("decode request body: %v", err)
	}
	return body
}

func assertNotExists(t *testing.T, err error) {
	t.Helper()

	var notExists *customErrors.NotExists
	if !errors.As(err, &notExists) {
		t.Errorf("error = %v (%T), want *errors.NotExists", err, err)
	}
}

func assertUnexpectedResult(t *testing.T, err error) {
	t.Helper()

	var unexpected *customErrors.UnexpectedResult
	if !errors.As(err, &unexpected) {
		t.Errorf("error = %v (%T), want *errors.UnexpectedResult", err, err)
	}
}

func Test_PreflightCheck(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		token string
	}{
		{name: "missing url", url: "", token: testToken},
		{name: "missing token", url: "https://authentik.example.com", token: ""},
		{name: "url without scheme", url: "authentik.example.com", token: testToken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.url, tt.token).GetFlowBySlug("default-authentication-flow")
			if err == nil {
				t.Fatal("GetFlowBySlug() error = nil, want preflight error")
			}
		})
	}
}

func Test_New_URLWithPathPrefix(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	r := New(srv.URL+"/authentik/", testToken)
	if err := r.DeleteGroup("group-uuid"); err != nil {
		t.Fatalf("DeleteGroup() error = %v", err)
	}

	if want := "/authentik/api/v3/core/groups/group-uuid/"; gotPath != want {
		t.Errorf("request path = %q, want %q", gotPath, want)
	}
}

func Test_ErrorBodyIsSurfaced(t *testing.T) {
	r := newFakeAuthentik(t, map[string]route{
		"DELETE /api/v3/core/groups/group-uuid/": respond(http.StatusForbidden, errorBody("no permissions")),
	})

	err := r.DeleteGroup("group-uuid")
	assertUnexpectedResult(t, err)
	if err == nil || !strings.Contains(err.Error(), "no permissions") {
		t.Errorf("error = %v, want it to contain the response detail", err)
	}
}

func Test_UndecodableResponseIsReported(t *testing.T) {
	// A role response without "pk", e.g. from a server whose API version does
	// not match the client.
	r := newFakeAuthentik(t, map[string]route{
		"POST /api/v3/rbac/roles/": respond(http.StatusCreated, obj{"name": "example-role"}),
	})

	_, err := r.CreateRole("example-role")
	assertUnexpectedResult(t, err)
	if err == nil || !strings.Contains(err.Error(), "pk") {
		t.Errorf("error = %v, want it to name the missing property", err)
	}
}
