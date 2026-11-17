package authentik

import (
	api "goauthentik.io/api/v3"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

// Keys of the custom attributes the CLI stores on authentik users and groups.
const (
	attrTenant   = "tenant"
	attrUserType = "userType"
)

func groupAttributesToAPI(attrs ak.GroupAttributes) map[string]interface{} {
	m := map[string]interface{}{}
	if attrs.Tenant != "" {
		m[attrTenant] = attrs.Tenant
	}
	return m
}

func userAttributesToAPI(attrs ak.UserAttributes) map[string]interface{} {
	m := map[string]interface{}{}
	if attrs.UserType != "" {
		m[attrUserType] = attrs.UserType
	}
	if attrs.Tenant != "" {
		m[attrTenant] = attrs.Tenant
	}
	return m
}

// stringAttr returns attrs[key] when it holds a string, or "" otherwise.
func stringAttr(attrs map[string]interface{}, key string) string {
	s, _ := attrs[key].(string)
	return s
}

func mapRole(role *api.Role) *ak.Role {
	return &ak.Role{
		PK:   role.Pk,
		Name: role.Name,
	}
}

func mapGroup(group *api.Group) *ak.Group {
	return &ak.Group{
		PK:   group.Pk,
		Name: group.Name,
		GroupAttributes: ak.GroupAttributes{
			Tenant: stringAttr(group.Attributes, attrTenant),
		},
	}
}

func mapGroupWithRoles(group *api.Group) *ak.Group {
	g := mapGroup(group)
	g.Roles = make([]ak.Role, 0, len(group.RolesObj))
	for i := range group.RolesObj {
		g.Roles = append(g.Roles, *mapRole(&group.RolesObj[i]))
	}
	return g
}

func mapUser(user *api.User) *ak.User {
	groups := make([]ak.Group, 0, len(user.GroupsObj))
	for _, group := range user.GroupsObj {
		groups = append(groups, ak.Group{
			PK:   group.Pk,
			Name: group.Name,
		})
	}

	return &ak.User{
		PK:          int(user.Pk),
		Username:    user.Username,
		Name:        user.Name,
		Email:       user.GetEmail(),
		Path:        user.GetPath(),
		IsActive:    user.GetIsActive(),
		IsSuperuser: user.IsSuperuser,
		Attributes: ak.UserAttributes{
			UserType: stringAttr(user.Attributes, attrUserType),
			Tenant:   stringAttr(user.Attributes, attrTenant),
		},
		Groups: groups,
	}
}

func mapOidcProvider(provider *api.OAuth2Provider) *ak.OidcProvider {
	var redirectUris []ak.OidcRedirectUri
	for _, uri := range provider.RedirectUris {
		redirectUris = append(redirectUris, ak.OidcRedirectUri{
			MatchingMode: string(uri.MatchingMode),
			Url:          uri.Url,
		})
	}

	return &ak.OidcProvider{
		Provider: ak.Provider{
			PK:                 int(provider.Pk),
			Name:               provider.Name,
			AuthenticationFlow: provider.GetAuthenticationFlow(),
			AuthorizationFlow:  provider.AuthorizationFlow,
			InvalidationFlow:   provider.InvalidationFlow,
		},
		PropertyMappings:     provider.PropertyMappings,
		ClientType:           string(provider.GetClientType()),
		ClientId:             provider.GetClientId(),
		ClientSecret:         provider.GetClientSecret(),
		AccessCodeValidity:   provider.GetAccessCodeValidity(),
		AccessTokenValidity:  provider.GetAccessTokenValidity(),
		RefreshTokenValidity: provider.GetRefreshTokenValidity(),
		SigningKey:           provider.GetSigningKey(),
		EncryptionKey:        provider.GetEncryptionKey(),
		RedirectUris:         redirectUris,
		SubMode:              string(provider.GetSubMode()),
		IssuerMode:           string(provider.GetIssuerMode()),
	}
}

func mapApplication(app *api.Application) *ak.Application {
	// provider_obj is null for applications without a provider; GetProviderObj
	// then returns the zero Provider.
	provider := app.GetProviderObj()

	return &ak.Application{
		PK:           app.Pk,
		Name:         app.Name,
		Slug:         app.Slug,
		ProviderPK:   int(provider.Pk),
		ProviderName: provider.Name,
		ProviderType: providerType(provider.MetaModelName),
	}
}

func providerType(metaModelName string) string {
	switch metaModelName {
	case ProviderTypeMetaModelLDAP:
		return constants.ProviderTypeLDAP
	case ProviderTypeMetaModelOIDC:
		return constants.ProviderTypeOIDC
	case ProviderTypeMetaModelProxy:
		return constants.ProviderTypeProxy
	case ProviderTypeMetaModelRAC:
		return constants.ProviderTypeRAC
	case ProviderTypeMetaModelSAML:
		return constants.ProviderTypeSAML
	case ProviderTypeMetaModelSCIM:
		return constants.ProviderTypeSCIM
	default:
		return ""
	}
}

func mapFlow(flow *api.Flow) ak.Flow {
	return ak.Flow{
		PK:          flow.Pk,
		Name:        flow.Name,
		Slug:        flow.Slug,
		Designation: string(flow.Designation),
	}
}
