package provider

import (
	"errors"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	customErrors "github.com/svetlyopet/authentik-cli/internal/errors"
	mock_ak "github.com/svetlyopet/authentik-cli/mocks/ak"
)

var defaultFlows = map[string]ak.Flow{
	FlowSlugAuthentication:               {PK: "authn-pk", Slug: FlowSlugAuthentication, Designation: "authentication"},
	FlowSlugAuthorizationExplicitConsent: {PK: "explicit-pk", Slug: FlowSlugAuthorizationExplicitConsent, Designation: "authorization"},
	FlowSlugAuthorizationImplicitConsent: {PK: "implicit-pk", Slug: FlowSlugAuthorizationImplicitConsent, Designation: "authorization"},
	FlowSlugInvalidation:                 {PK: "invalidation-pk", Slug: FlowSlugInvalidation, Designation: "invalidation"},
}

// setRepo installs a mock repository as ak.Repo for the duration of the test.
func setRepo(t *testing.T) *mock_ak.MockAuthentikRepository {
	t.Helper()

	repo := mock_ak.NewMockAuthentikRepository(gomock.NewController(t))
	prev := ak.Repo
	ak.Repo = repo
	t.Cleanup(func() { ak.Repo = prev })

	return repo
}

// serveFlows makes the mock answer flow lookups from flows; slugs missing
// from flows are reported as not existing.
func serveFlows(repo *mock_ak.MockAuthentikRepository, flows map[string]ak.Flow) {
	repo.EXPECT().GetFlowBySlug(gomock.Any()).DoAndReturn(func(slug string) (*ak.Flow, error) {
		f, ok := flows[slug]
		if !ok {
			return nil, customErrors.NewNotExists("get flow " + slug + ": not found")
		}
		return &f, nil
	}).AnyTimes()
}

func flowsWithout(slug string) map[string]ak.Flow {
	flows := map[string]ak.Flow{}
	for k, v := range defaultFlows {
		if k != slug {
			flows[k] = v
		}
	}
	return flows
}

func createProvider(consentType string) (*ak.OidcProvider, error) {
	return CreateOidcProvider("example-app", OidcClientTypeConfidential, consentType, false,
		[]string{"http://localhost:8000"}, nil)
}

func Test_CreateOidcProvider_UsesDefaultFlows(t *testing.T) {
	tests := []struct {
		consentType       string
		wantAuthorization string
	}{
		{consentType: OidcConsentTypeExplicit, wantAuthorization: "explicit-pk"},
		{consentType: OidcConsentTypeImplicit, wantAuthorization: "implicit-pk"},
	}

	for _, tt := range tests {
		t.Run(tt.consentType, func(t *testing.T) {
			repo := setRepo(t)
			serveFlows(repo, defaultFlows)
			repo.EXPECT().CreateOidcProvider(gomock.Any()).DoAndReturn(func(p ak.OidcProvider) (*ak.OidcProvider, error) {
				if p.AuthenticationFlow != "authn-pk" {
					t.Errorf("AuthenticationFlow = %q, want authn-pk", p.AuthenticationFlow)
				}
				if p.AuthorizationFlow != tt.wantAuthorization {
					t.Errorf("AuthorizationFlow = %q, want %q", p.AuthorizationFlow, tt.wantAuthorization)
				}
				if p.InvalidationFlow != "invalidation-pk" {
					t.Errorf("InvalidationFlow = %q, want invalidation-pk", p.InvalidationFlow)
				}
				p.PK = 1
				return &p, nil
			})

			if _, err := createProvider(tt.consentType); err != nil {
				t.Fatalf("CreateOidcProvider() error = %v", err)
			}
		})
	}
}

func Test_CreateOidcProvider_OptionalAuthenticationFlow(t *testing.T) {
	repo := setRepo(t)
	serveFlows(repo, flowsWithout(FlowSlugAuthentication))
	repo.EXPECT().CreateOidcProvider(gomock.Any()).DoAndReturn(func(p ak.OidcProvider) (*ak.OidcProvider, error) {
		if p.AuthenticationFlow != "" {
			t.Errorf("AuthenticationFlow = %q, want empty", p.AuthenticationFlow)
		}
		return &p, nil
	})

	if _, err := createProvider(OidcConsentTypeExplicit); err != nil {
		t.Fatalf("CreateOidcProvider() error = %v", err)
	}
}

func Test_CreateOidcProvider_MissingRequiredFlow(t *testing.T) {
	for _, slug := range []string{FlowSlugAuthorizationExplicitConsent, FlowSlugInvalidation} {
		t.Run(slug, func(t *testing.T) {
			repo := setRepo(t)
			serveFlows(repo, flowsWithout(slug))
			// No CreateOidcProvider expectation: the mock fails the test if it is called.

			_, err := createProvider(OidcConsentTypeExplicit)

			var notExists *customErrors.NotExists
			if !errors.As(err, &notExists) {
				t.Fatalf("error = %v (%T), want *errors.NotExists", err, err)
			}
			if !strings.Contains(err.Error(), slug) {
				t.Errorf("error = %q, want it to name %s", err, slug)
			}
		})
	}
}

func Test_CreateOidcProvider_InvalidFlow(t *testing.T) {
	tests := []struct {
		name string
		flow ak.Flow
		want string
	}{
		{
			name: "wrong designation",
			flow: ak.Flow{PK: "pk", Slug: FlowSlugInvalidation, Designation: "authentication"},
			want: "designation",
		},
		{
			name: "missing primary key",
			flow: ak.Flow{Slug: FlowSlugInvalidation, Designation: "invalidation"},
			want: "primary key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flows := flowsWithout(FlowSlugInvalidation)
			flows[FlowSlugInvalidation] = tt.flow
			repo := setRepo(t)
			serveFlows(repo, flows)

			_, err := createProvider(OidcConsentTypeExplicit)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want it to mention %q", err, tt.want)
			}
		})
	}
}

func Test_CreateOidcProvider_LookupError(t *testing.T) {
	repo := setRepo(t)
	repo.EXPECT().GetFlowBySlug(gomock.Any()).Return(nil, errors.New("get flow: decoding response: boom"))

	_, err := createProvider(OidcConsentTypeExplicit)
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("error = %v, want the lookup error", err)
	}
}

func Test_CreateOidcProvider_InvalidConsentType(t *testing.T) {
	setRepo(t) // no calls expected

	if _, err := createProvider("sometimes"); err == nil {
		t.Fatal("CreateOidcProvider() error = nil, want error")
	}
}
