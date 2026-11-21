package authentik

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/svetlyopet/authentik-cli/internal/ak"
)

func Test_CreateGroup(t *testing.T) {
	tests := []struct {
		name    string
		route   route
		want    *ak.Group
		wantErr bool
	}{
		{
			name:    "create group error",
			route:   respond(http.StatusBadRequest, errorBody("bad request")),
			wantErr: true,
		},
		{
			name: "create group success",
			route: func(t *testing.T, r *http.Request) (int, any) {
				body := decodeBody(t, r)
				if body["name"] != "example-group" || body["is_superuser"] != false {
					t.Errorf("unexpected request body %v", body)
				}
				if roles, _ := body["roles"].([]any); len(roles) != 1 || roles[0] != "example-role" {
					t.Errorf("roles = %v, want [example-role]", body["roles"])
				}
				if attrs, _ := body["attributes"].(map[string]any); attrs["tenant"] != "example-tenant" {
					t.Errorf("attributes = %v, want tenant example-tenant", body["attributes"])
				}
				return http.StatusCreated, groupFixture("random-uuid", "example-group", obj{
					"attributes": obj{"tenant": "example-tenant"},
				})
			},
			want: &ak.Group{
				PK:              "random-uuid",
				Name:            "example-group",
				GroupAttributes: ak.GroupAttributes{Tenant: "example-tenant"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeAuthentik(t, map[string]route{"POST /api/v3/core/groups/": tt.route})

			got, err := r.CreateGroup("example-group", []string{"example-role"}, ak.GroupAttributes{Tenant: "example-tenant"})
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateGroup() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_GetGroupByName(t *testing.T) {
	group := groupFixture("random-uuid", "example-group", obj{"attributes": obj{"tenant": "example-tenant"}})

	tests := []struct {
		name          string
		route         route
		want          *ak.Group
		wantNotExists bool
		wantErr       bool
	}{
		{
			name:          "group not found",
			route:         respond(http.StatusOK, paginated(0)),
			wantErr:       true,
			wantNotExists: true,
		},
		{
			name:    "multiple groups found",
			route:   respond(http.StatusOK, paginated(0, group, group)),
			wantErr: true,
		},
		{
			name:    "api error",
			route:   respond(http.StatusForbidden, errorBody("no permissions")),
			wantErr: true,
		},
		{
			name: "group found",
			route: func(t *testing.T, r *http.Request) (int, any) {
				if got := r.URL.Query().Get("name"); got != "example-group" {
					t.Errorf("name query = %q, want example-group", got)
				}
				return http.StatusOK, paginated(0, group)
			},
			want: &ak.Group{
				PK:              "random-uuid",
				Name:            "example-group",
				GroupAttributes: ak.GroupAttributes{Tenant: "example-tenant"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeAuthentik(t, map[string]route{"GET /api/v3/core/groups/": tt.route})

			got, err := r.GetGroupByName("example-group")
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetGroupByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNotExists {
				assertNotExists(t, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGroupByName() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_GetGroup(t *testing.T) {
	t.Run("group with roles", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"GET /api/v3/core/groups/random-uuid/": respond(http.StatusOK, groupFixture("random-uuid", "example-group", obj{
				"roles_obj": list(roleFixture("role-uuid", "example-role")),
			})),
		})

		got, err := r.GetGroup("random-uuid")
		if err != nil {
			t.Fatalf("GetGroup() error = %v", err)
		}
		want := &ak.Group{
			PK:    "random-uuid",
			Name:  "example-group",
			Roles: []ak.Role{{PK: "role-uuid", Name: "example-role"}},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetGroup() got = %+v, want %+v", got, want)
		}
	})

	t.Run("group not found", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"GET /api/v3/core/groups/random-uuid/": respond(http.StatusNotFound, errorBody("not found")),
		})

		_, err := r.GetGroup("random-uuid")
		assertNotExists(t, err)
	})
}

func Test_DeleteGroup(t *testing.T) {
	tests := []struct {
		name          string
		route         route
		wantErr       bool
		wantNotExists bool
	}{
		{name: "deleted", route: respond(http.StatusNoContent, nil)},
		{name: "not found", route: respond(http.StatusNotFound, errorBody("not found")), wantErr: true, wantNotExists: true},
		{name: "forbidden", route: respond(http.StatusForbidden, errorBody("no permissions")), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeAuthentik(t, map[string]route{"DELETE /api/v3/core/groups/random-uuid/": tt.route})

			err := r.DeleteGroup("random-uuid")
			if (err != nil) != tt.wantErr {
				t.Fatalf("DeleteGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNotExists {
				assertNotExists(t, err)
			}
		})
	}
}

func Test_CreateUser(t *testing.T) {
	user := ak.User{
		Username: "johndoe",
		Name:     "John Doe",
		Email:    "John.Doe@example.com",
		Path:     "users",
		IsActive: true,
		Attributes: ak.UserAttributes{
			UserType: "tenantAdmin",
			Tenant:   "example-tenant",
		},
	}

	tests := []struct {
		name    string
		route   route
		want    *ak.User
		wantErr bool
	}{
		{
			name:    "create user token error",
			route:   respond(http.StatusForbidden, errorBody("no permissions")),
			wantErr: true,
		},
		{
			name: "create user payload error",
			route: respond(http.StatusBadRequest, map[string][]string{
				"email": {"Enter a valid email address."},
			}),
			wantErr: true,
		},
		{
			name: "create user success",
			route: func(t *testing.T, r *http.Request) (int, any) {
				body := decodeBody(t, r)
				for key, want := range map[string]any{
					"username":  "johndoe",
					"name":      "John Doe",
					"email":     "John.Doe@example.com",
					"path":      "users",
					"is_active": true,
				} {
					if body[key] != want {
						t.Errorf("request %s = %v, want %v", key, body[key], want)
					}
				}
				attrs, _ := body["attributes"].(map[string]any)
				if attrs["userType"] != "tenantAdmin" || attrs["tenant"] != "example-tenant" {
					t.Errorf("request attributes = %v", attrs)
				}

				return http.StatusCreated, userFixture(999, "johndoe", "John Doe", obj{
					"email":      "John.Doe@example.com",
					"path":       "users",
					"is_active":  true,
					"attributes": obj{"userType": "tenantAdmin", "tenant": "example-tenant"},
				})
			},
			want: &ak.User{
				PK:       999,
				Username: "johndoe",
				Name:     "John Doe",
				Email:    "John.Doe@example.com",
				Path:     "users",
				IsActive: true,
				Attributes: ak.UserAttributes{
					UserType: "tenantAdmin",
					Tenant:   "example-tenant",
				},
				Groups: []ak.Group{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeAuthentik(t, map[string]route{"POST /api/v3/core/users/": tt.route})

			got, err := r.CreateUser(user)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_GetUserByUsername(t *testing.T) {
	user := userFixture(7, "johndoe", "John Doe", obj{
		"groups_obj": list(partialGroupFixture("group-uuid", "example-group")),
	})

	tests := []struct {
		name          string
		route         route
		want          *ak.User
		wantErr       bool
		wantNotExists bool
	}{
		{
			name:          "user not found",
			route:         respond(http.StatusOK, paginated(0)),
			wantErr:       true,
			wantNotExists: true,
		},
		{
			name:    "multiple users found",
			route:   respond(http.StatusOK, paginated(0, user, user)),
			wantErr: true,
		},
		{
			name: "user found",
			route: func(t *testing.T, r *http.Request) (int, any) {
				if got := r.URL.Query().Get("username"); got != "johndoe" {
					t.Errorf("username query = %q, want johndoe", got)
				}
				return http.StatusOK, paginated(0, user)
			},
			want: &ak.User{
				PK:       7,
				Username: "johndoe",
				Name:     "John Doe",
				Groups:   []ak.Group{{PK: "group-uuid", Name: "example-group"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeAuthentik(t, map[string]route{"GET /api/v3/core/users/": tt.route})

			got, err := r.GetUserByUsername("johndoe")
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNotExists {
				assertNotExists(t, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByUsername() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_DeleteUser(t *testing.T) {
	t.Run("deleted", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{"DELETE /api/v3/core/users/7/": respond(http.StatusNoContent, nil)})
		if err := r.DeleteUser("7"); err != nil {
			t.Errorf("DeleteUser() error = %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		r := newFakeAuthentik(t, nil)
		if err := r.DeleteUser("not-a-number"); err == nil {
			t.Error("DeleteUser() error = nil, want error")
		}
	})
}

func Test_AddUserToGroup(t *testing.T) {
	r := newFakeAuthentik(t, map[string]route{
		"POST /api/v3/core/groups/group-uuid/add_user/": func(t *testing.T, r *http.Request) (int, any) {
			if body := decodeBody(t, r); body["pk"] != float64(7) {
				t.Errorf("request pk = %v, want 7", body["pk"])
			}
			return http.StatusNoContent, nil
		},
	})

	if err := r.AddUserToGroup(7, "group-uuid"); err != nil {
		t.Errorf("AddUserToGroup() error = %v", err)
	}
}

func Test_CreateApplication(t *testing.T) {
	r := newFakeAuthentik(t, map[string]route{
		"POST /api/v3/core/applications/": func(t *testing.T, r *http.Request) (int, any) {
			body := decodeBody(t, r)
			if body["name"] != "example-app" || body["slug"] != "example-app" || body["provider"] != float64(3) {
				t.Errorf("unexpected request body %v", body)
			}
			return http.StatusCreated, applicationFixture("app-uuid", "example-app", "example-app",
				providerFixture(3, "example-provider", ProviderTypeMetaModelOIDC))
		},
	})

	got, err := r.CreateApplication("example-app", "example-app", 3)
	if err != nil {
		t.Fatalf("CreateApplication() error = %v", err)
	}
	want := &ak.Application{
		PK:           "app-uuid",
		Name:         "example-app",
		Slug:         "example-app",
		ProviderPK:   3,
		ProviderName: "example-provider",
		ProviderType: "oidc",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateApplication() got = %+v, want %+v", got, want)
	}
}

func Test_GetApplicationByName(t *testing.T) {
	app := applicationFixture("app-uuid", "example-app", "example-app",
		providerFixture(3, "example-provider", ProviderTypeMetaModelOIDC))

	t.Run("found", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"GET /api/v3/core/applications/": respond(http.StatusOK, paginated(0, app)),
		})

		got, err := r.GetApplicationByName("example-app")
		if err != nil {
			t.Fatalf("GetApplicationByName() error = %v", err)
		}
		if got.PK != "app-uuid" || got.ProviderPK != 3 || got.ProviderType != "oidc" {
			t.Errorf("GetApplicationByName() got = %+v", got)
		}
	})

	t.Run("not found", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"GET /api/v3/core/applications/": respond(http.StatusOK, paginated(0)),
		})

		_, err := r.GetApplicationByName("example-app")
		assertNotExists(t, err)
	})
}

func Test_DeleteApplication(t *testing.T) {
	t.Run("deleted", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"DELETE /api/v3/core/applications/example-app/": respond(http.StatusNoContent, nil),
		})
		if err := r.DeleteApplication("example-app"); err != nil {
			t.Errorf("DeleteApplication() error = %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"DELETE /api/v3/core/applications/example-app/": respond(http.StatusNotFound, errorBody("not found")),
		})
		assertNotExists(t, r.DeleteApplication("example-app"))
	})
}

func Test_GetApplicationByName_WithoutProvider(t *testing.T) {
	r := newFakeAuthentik(t, map[string]route{
		"GET /api/v3/core/applications/": respond(http.StatusOK, paginated(0,
			applicationFixture("app-uuid", "example-app", "example-app", nil))),
	})

	got, err := r.GetApplicationByName("example-app")
	if err != nil {
		t.Fatalf("GetApplicationByName() error = %v", err)
	}
	want := &ak.Application{PK: "app-uuid", Name: "example-app", Slug: "example-app"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetApplicationByName() got = %+v, want %+v", got, want)
	}
}

// assertBody fails the test unless the JSON request body equals want exactly,
// which proves that a partial update sends no field it was not asked to.
func assertBody(t *testing.T, r *http.Request, want map[string]any) {
	t.Helper()

	if got := decodeBody(t, r); !reflect.DeepEqual(got, want) {
		t.Errorf("request body = %v, want exactly %v", got, want)
	}
}

func Test_UpdateUser(t *testing.T) {
	tests := []struct {
		name     string
		patch    ak.UserPatch
		wantBody map[string]any
	}{
		{
			name:     "single field",
			patch:    ak.UserPatch{Email: new("new@example.com")},
			wantBody: map[string]any{"email": "new@example.com"},
		},
		{
			name:     "deactivate",
			patch:    ak.UserPatch{IsActive: new(false)},
			wantBody: map[string]any{"is_active": false},
		},
		{
			name:  "several fields",
			patch: ak.UserPatch{Username: new("janedoe"), Name: new("Jane Doe")},
			wantBody: map[string]any{
				"username": "janedoe",
				"name":     "Jane Doe",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeAuthentik(t, map[string]route{
				"PATCH /api/v3/core/users/7/": func(t *testing.T, r *http.Request) (int, any) {
					assertBody(t, r, tt.wantBody)
					return http.StatusOK, userFixture(7, "johndoe", "John Doe", nil)
				},
			})

			got, err := r.UpdateUser(7, tt.patch)
			if err != nil {
				t.Fatalf("UpdateUser() error = %v", err)
			}
			if got.PK != 7 {
				t.Errorf("UpdateUser() PK = %d, want 7", got.PK)
			}
		})
	}

	t.Run("not found", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"PATCH /api/v3/core/users/7/": respond(http.StatusNotFound, errorBody("not found")),
		})

		_, err := r.UpdateUser(7, ak.UserPatch{Name: new("x")})
		assertNotExists(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		r := newFakeAuthentik(t, map[string]route{
			"PATCH /api/v3/core/users/7/": respond(http.StatusBadRequest, map[string][]string{
				"email": {"Enter a valid email address."},
			}),
		})

		if _, err := r.UpdateUser(7, ak.UserPatch{Email: new("nope")}); err == nil {
			t.Error("UpdateUser() error = nil, want error")
		}
	})
}

func Test_UpdateGroup(t *testing.T) {
	r := newFakeAuthentik(t, map[string]route{
		"PATCH /api/v3/core/groups/group-uuid/": func(t *testing.T, r *http.Request) (int, any) {
			assertBody(t, r, map[string]any{"name": "renamed-group"})
			return http.StatusOK, groupFixture("group-uuid", "renamed-group", nil)
		},
	})

	got, err := r.UpdateGroup("group-uuid", ak.GroupPatch{Name: new("renamed-group")})
	if err != nil {
		t.Fatalf("UpdateGroup() error = %v", err)
	}
	want := &ak.Group{PK: "group-uuid", Name: "renamed-group"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateGroup() got = %+v, want %+v", got, want)
	}
}

func Test_UpdateApplication(t *testing.T) {
	r := newFakeAuthentik(t, map[string]route{
		"PATCH /api/v3/core/applications/example-app/": func(t *testing.T, r *http.Request) (int, any) {
			assertBody(t, r, map[string]any{"slug": "new-slug"})
			return http.StatusOK, applicationFixture("app-uuid", "example-app", "new-slug",
				providerFixture(3, "example-app", ProviderTypeMetaModelOIDC))
		},
	})

	got, err := r.UpdateApplication("example-app", ak.ApplicationPatch{Slug: new("new-slug")})
	if err != nil {
		t.Fatalf("UpdateApplication() error = %v", err)
	}
	want := &ak.Application{
		PK:           "app-uuid",
		Name:         "example-app",
		Slug:         "new-slug",
		ProviderPK:   3,
		ProviderName: "example-app",
		ProviderType: "oidc",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateApplication() got = %+v, want %+v", got, want)
	}
}
