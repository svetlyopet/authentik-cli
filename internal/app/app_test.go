package app

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/provider"
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

func oidcApplication(providerName string) *ak.Application {
	return &ak.Application{
		PK:           "app-uuid",
		Name:         "example-app",
		Slug:         "example-app",
		ProviderPK:   3,
		ProviderType: constants.ProviderTypeOIDC,
		ProviderName: providerName,
	}
}

func oidcProvider(name string) *ak.OidcProvider {
	return &ak.OidcProvider{
		Provider:   ak.Provider{PK: 3, Name: name},
		ClientType: provider.OidcClientTypeConfidential,
	}
}

func Test_Update_RenameAlsoRenamesMatchingProvider(t *testing.T) {
	repo := setRepo(t)
	repo.EXPECT().GetApplicationByName("example-app").Return(oidcApplication("example-app"), nil)
	repo.EXPECT().GetOidcProvider(3).Return(oidcProvider("example-app"), nil)
	repo.EXPECT().UpdateOidcProvider(3, ak.OidcProviderPatch{Name: new("renamed")}).
		Return(oidcProvider("renamed"), nil)
	repo.EXPECT().UpdateApplication("example-app", ak.ApplicationPatch{Name: new("renamed")}).
		Return(&ak.Application{Name: "renamed"}, nil)

	if err := Update("example-app", Changes{Name: new("renamed")}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
}

func Test_Update_RenameKeepsCustomProviderName(t *testing.T) {
	repo := setRepo(t)
	repo.EXPECT().GetApplicationByName("example-app").Return(oidcApplication("custom-provider"), nil)
	repo.EXPECT().UpdateApplication("example-app", ak.ApplicationPatch{Name: new("renamed")}).
		Return(&ak.Application{Name: "renamed"}, nil)

	if err := Update("example-app", Changes{Name: new("renamed")}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
}

func Test_Update_ProviderOnly(t *testing.T) {
	repo := setRepo(t)
	repo.EXPECT().GetApplicationByName("example-app").Return(oidcApplication("example-app"), nil)
	repo.EXPECT().GetOidcProvider(3).Return(oidcProvider("example-app"), nil)
	repo.EXPECT().UpdateOidcProvider(3, ak.OidcProviderPatch{ClientType: new(provider.OidcClientTypePublic)}).
		Return(oidcProvider("example-app"), nil)
	// No UpdateApplication expectation: the application itself is not touched.

	changes := Changes{Oidc: provider.OidcChanges{ClientType: new(provider.OidcClientTypePublic)}}
	if err := Update("example-app", changes); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
}

func Test_Update_OidcChangesOnNonOidcApplication(t *testing.T) {
	repo := setRepo(t)
	repo.EXPECT().GetApplicationByName("example-app").Return(&ak.Application{
		Name:         "example-app",
		Slug:         "example-app",
		ProviderType: constants.ProviderTypeSAML,
	}, nil)

	changes := Changes{Oidc: provider.OidcChanges{ClientType: new(provider.OidcClientTypePublic)}}
	if err := Update("example-app", changes); err == nil {
		t.Fatal("Update() error = nil, want error")
	}
}
