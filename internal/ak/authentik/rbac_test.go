package authentik

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/svetlyopet/authentik-cli/internal/ak"
)

func Test_CreateRole(t *testing.T) {
	tests := []struct {
		name    string
		route   route
		want    *ak.Role
		wantErr bool
	}{
		{
			name:    "create role error",
			route:   respond(http.StatusBadRequest, errorBody("bad request")),
			wantErr: true,
		},
		{
			name: "create role success",
			route: func(t *testing.T, r *http.Request) (int, any) {
				if body := decodeBody(t, r); body["name"] != "example-role" {
					t.Errorf("request name = %v, want example-role", body["name"])
				}
				return http.StatusCreated, roleFixture("random-uuid", "example-role")
			},
			want: &ak.Role{PK: "random-uuid", Name: "example-role"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeAuthentik(t, map[string]route{"POST /api/v3/rbac/roles/": tt.route})

			got, err := r.CreateRole("example-role")
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRole() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_GetRoleByName(t *testing.T) {
	roles := func(names ...string) obj {
		var results []obj
		for i, name := range names {
			results = append(results, roleFixture(fmt.Sprintf("uuid-%d", i+1), name))
		}
		return paginated(0, results...)
	}

	tests := []struct {
		name          string
		role          string
		route         route
		want          *ak.Role
		wantErr       bool
		wantNotExists bool
	}{
		{
			name:          "role not found",
			role:          "missing-role",
			route:         respond(http.StatusOK, roles()),
			wantErr:       true,
			wantNotExists: true,
		},
		{
			name:          "only fuzzy matches",
			role:          "example-role",
			route:         respond(http.StatusOK, roles("example-role-2", "other-example-role")),
			wantErr:       true,
			wantNotExists: true,
		},
		{
			name:    "get role by name error",
			role:    "example-role",
			route:   respond(http.StatusBadRequest, errorBody("bad request")),
			wantErr: true,
		},
		{
			name:    "multiple roles found",
			role:    "duplicate-role",
			route:   respond(http.StatusOK, roles("duplicate-role", "duplicate-role")),
			wantErr: true,
		},
		{
			name: "get role by name success",
			role: "example-role",
			route: func(t *testing.T, r *http.Request) (int, any) {
				if got := r.URL.Query().Get("search"); got != "example-role" {
					t.Errorf("search query = %q, want example-role", got)
				}
				return http.StatusOK, roles("example-role-2", "example-role")
			},
			want: &ak.Role{PK: "uuid-2", Name: "example-role"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeAuthentik(t, map[string]route{"GET /api/v3/rbac/roles/": tt.route})

			got, err := r.GetRoleByName(tt.role)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetRoleByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNotExists {
				assertNotExists(t, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRoleByName() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_DeleteRole(t *testing.T) {
	tests := []struct {
		name    string
		route   route
		wantErr bool
	}{
		{name: "delete role error", route: respond(http.StatusBadRequest, errorBody("bad request")), wantErr: true},
		{name: "delete role success", route: respond(http.StatusNoContent, nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeAuthentik(t, map[string]route{"DELETE /api/v3/rbac/roles/random-uuid/": tt.route})

			if err := r.DeleteRole("random-uuid"); (err != nil) != tt.wantErr {
				t.Errorf("DeleteRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_AssignViewPermissionsToTenantRole(t *testing.T) {
	const path = "POST /api/v3/rbac/permissions/assigned_by_roles/random-uuid/assign/"

	tests := []struct {
		name    string
		route   route
		wantErr bool
	}{
		{
			name:    "assign permissions error",
			route:   respond(http.StatusBadRequest, errorBody("bad request")),
			wantErr: true,
		},
		{
			name: "assign permissions success",
			route: func(t *testing.T, r *http.Request) (int, any) {
				perms, _ := decodeBody(t, r)["permissions"].([]any)
				if len(perms) != len(tenantRoleViewPermissions) {
					t.Errorf("permissions = %v, want %v", perms, tenantRoleViewPermissions)
				}
				return http.StatusOK, list()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeAuthentik(t, map[string]route{path: tt.route})

			if err := r.AssignViewPermissionsToTenantRole("random-uuid"); (err != nil) != tt.wantErr {
				t.Errorf("AssignViewPermissionsToTenantRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
