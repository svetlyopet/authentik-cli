package core

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	mock_ak "github.com/svetlyopet/authentik-cli/mocks/ak"
)

// setRepo installs a mock repository as ak.Repo for the duration of the test.
// Calls without a matching expectation fail the test.
func setRepo(t *testing.T) *mock_ak.MockAuthentikRepository {
	t.Helper()

	repo := mock_ak.NewMockAuthentikRepository(gomock.NewController(t))
	prev := ak.Repo
	ak.Repo = repo
	t.Cleanup(func() { ak.Repo = prev })

	return repo
}

func Test_UpdateUser(t *testing.T) {
	current := ak.User{
		PK:       7,
		Username: "johndoe",
		Name:     "John Doe",
		Email:    "john@example.com",
		IsActive: true,
	}

	tests := []struct {
		name    string
		changes UserChanges
		// wantPatch is the patch sent to authentik, or nil when no update
		// call is expected.
		wantPatch *ak.UserPatch
	}{
		{
			name:      "single field",
			changes:   UserChanges{Email: new("jd@example.com")},
			wantPatch: &ak.UserPatch{Email: new("jd@example.com")},
		},
		{
			name:      "name and surname form the display name",
			changes:   UserChanges{Name: new("Johnny"), Surname: new("Doe")},
			wantPatch: &ak.UserPatch{Name: new("Johnny Doe")},
		},
		{
			name:      "name without surname",
			changes:   UserChanges{Name: new("Johnny")},
			wantPatch: &ak.UserPatch{Name: new("Johnny")},
		},
		{
			name:      "deactivate",
			changes:   UserChanges{IsActive: new(false)},
			wantPatch: &ak.UserPatch{IsActive: new(false)},
		},
		{
			name:      "values equal to the current ones are not sent",
			changes:   UserChanges{Username: new("johndoe"), Email: new("new@example.com")},
			wantPatch: &ak.UserPatch{Email: new("new@example.com")},
		},
		{
			name: "nothing changes",
			changes: UserChanges{
				Name:     new("John"),
				Surname:  new("Doe"),
				Email:    new("john@example.com"),
				IsActive: new(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setRepo(t)
			user := current
			repo.EXPECT().GetUserByUsername("johndoe").Return(&user, nil)
			if tt.wantPatch != nil {
				repo.EXPECT().UpdateUser(7, *tt.wantPatch).Return(&user, nil)
			}

			if err := UpdateUser("johndoe", tt.changes); err != nil {
				t.Fatalf("UpdateUser() error = %v", err)
			}
		})
	}

	t.Run("surname without name", func(t *testing.T) {
		setRepo(t) // no calls expected

		if err := UpdateUser("johndoe", UserChanges{Surname: new("Doe")}); err == nil {
			t.Fatal("UpdateUser() error = nil, want error")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		repo := setRepo(t)
		repo.EXPECT().GetUserByUsername("johndoe").Return(nil, customErrors.NewNotExists("user not found"))

		if err := UpdateUser("johndoe", UserChanges{Email: new("x@example.com")}); err == nil {
			t.Fatal("UpdateUser() error = nil, want error")
		}
	})
}

func Test_UpdateGroup(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		repo := setRepo(t)
		repo.EXPECT().GetGroupByName("old-name").Return(&ak.Group{PK: "group-uuid", Name: "old-name"}, nil)
		repo.EXPECT().UpdateGroup("group-uuid", ak.GroupPatch{Name: new("new-name")}).
			Return(&ak.Group{PK: "group-uuid", Name: "new-name"}, nil)

		if err := UpdateGroup("old-name", GroupChanges{Name: new("new-name")}); err != nil {
			t.Fatalf("UpdateGroup() error = %v", err)
		}
	})

	t.Run("same name", func(t *testing.T) {
		repo := setRepo(t)
		repo.EXPECT().GetGroupByName("old-name").Return(&ak.Group{PK: "group-uuid", Name: "old-name"}, nil)

		if err := UpdateGroup("old-name", GroupChanges{Name: new("old-name")}); err != nil {
			t.Fatalf("UpdateGroup() error = %v", err)
		}
	})

	t.Run("tenant group cannot be renamed", func(t *testing.T) {
		repo := setRepo(t)
		repo.EXPECT().GetGroupByName("acme-tenant-admin").Return(&ak.Group{
			PK:              "group-uuid",
			Name:            "acme-tenant-admin",
			GroupAttributes: ak.GroupAttributes{Tenant: "acme"},
		}, nil)

		if err := UpdateGroup("acme-tenant-admin", GroupChanges{Name: new("other")}); err == nil {
			t.Fatal("UpdateGroup() error = nil, want error")
		}
	})
}

func Test_UpdateApplication(t *testing.T) {
	application := &ak.Application{PK: "app-uuid", Name: "example-app", Slug: "example-app"}

	t.Run("only changed fields are sent", func(t *testing.T) {
		repo := setRepo(t)
		repo.EXPECT().UpdateApplication("example-app", ak.ApplicationPatch{Slug: new("new-slug")}).
			Return(&ak.Application{Name: "example-app", Slug: "new-slug"}, nil)

		err := UpdateApplication(application, ApplicationChanges{Name: new("example-app"), Slug: new("new-slug")})
		if err != nil {
			t.Fatalf("UpdateApplication() error = %v", err)
		}
	})

	t.Run("nothing changes", func(t *testing.T) {
		setRepo(t) // no calls expected

		if err := UpdateApplication(application, ApplicationChanges{Slug: new("example-app")}); err != nil {
			t.Fatalf("UpdateApplication() error = %v", err)
		}
	})
}
