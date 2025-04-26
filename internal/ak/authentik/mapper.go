package authentik

import (
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

func mapToCreateOrUpdateRoleResponse(role *getRoleResponse) *ak.Role {
	return &ak.Role{
		PK:   role.PK,
		Name: role.Name,
	}
}

func mapToGetRoleByNameResponse(roles *getRolesResponse) *ak.Role {
	res := &ak.Role{}
	for _, role := range roles.Results {
		res = &ak.Role{
			PK:   role.PK,
			Name: role.Name,
		}
	}

	return res
}

func mapToCreateOrUpdateGroupResponse(group *getGroupResponse) *ak.Group {
	return &ak.Group{
		PK:   group.PK,
		Name: group.Name,
		GroupAttributes: ak.GroupAttributes{
			Tenant: group.Attributes.Tenant,
		},
	}
}

func mapToGetGroupResponse(group *getGroupResponse) *ak.Group {
	var roles = []ak.Role{}

	for _, role := range group.RolesObj {
		roles = append(roles, ak.Role{
			PK:   role.PK,
			Name: role.Name,
		})
	}

	return &ak.Group{
		PK:   group.PK,
		Name: group.Name,
		GroupAttributes: ak.GroupAttributes{
			Tenant: group.Attributes.Tenant,
		},
		Roles: roles,
	}
}

func mapToGetGroupByNameResponse(groups *getGroupsResponse) *ak.Group {
	res := &ak.Group{}
	for _, group := range groups.Results {
		res = &ak.Group{
			PK:   group.PK,
			Name: group.Name,
			GroupAttributes: ak.GroupAttributes{
				Tenant: group.Attributes.Tenant,
			},
		}
	}

	return res
}

func mapToCreateOrUpdateUserResponse(user *getUserResponse) *ak.User {
	return &ak.User{
		PK:       user.PK,
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Path:     user.Path,
		IsActive: user.IsActive,
		Attributes: ak.UserAttributes{
			UserType: user.Attributes.UserType,
			Tenant:   user.Attributes.Tenant,
		},
	}
}

func mapToUsersGetResponse(users *getUsersResponse) *ak.User {
	var userGetResponse ak.User

	for _, userResults := range users.Results {
		userGetResponse.PK = userResults.PK
		userGetResponse.Username = userResults.Username
		userGetResponse.Name = userResults.Name
		userGetResponse.Email = userResults.Email
		userGetResponse.Path = userResults.Path
		userGetResponse.IsActive = userResults.IsActive
		userGetResponse.Attributes.UserType = userResults.Attributes.UserType
	}

	return &userGetResponse
}

func mapToGetOidcProviderResponse(provider *getOidcProviderResponse) *ak.OidcProvider {
	var redirectUris []ak.OidcRedirectUri

	for _, redirectUri := range provider.RedirectUris {
		redirectUris = append(redirectUris, ak.OidcRedirectUri{
			MatchingMode: redirectUri.MatchingMode,
			Url:          redirectUri.Url,
		})
	}

	providerResp := ak.OidcProvider{
		Provider: ak.Provider{
			PK:                 provider.PK,
			Name:               provider.Name,
			AuthenticationFlow: provider.AuthenticationFlow,
			AuthorizationFlow:  provider.AuthorizationFlow,
			InvalidationFlow:   provider.IssuerMode,
		},
		PropertyMappings: provider.PropertyMappings,
		ClientType:       provider.ClientType,
		ClientId:         provider.ClientId,
		ClientSecret:     provider.ClientSecret,
		SigningKey:       provider.SigningKey,
		EncryptionKey:    provider.EncryptionKey,
		RedirectUris:     redirectUris,
		SubMode:          provider.SubMode,
		IssuerMode:       provider.IssuerMode,
	}

	return &providerResp
}

func mapToCreateOrUpdateApplicationResponse(application *getApplicationResponse) *ak.Application {
	applicationResp := ak.Application{
		PK:           application.PK,
		Name:         application.Name,
		Slug:         application.Slug,
		ProviderPK:   application.ProviderObj.PK,
		ProviderName: application.ProviderObj.Name,
	}

	return &applicationResp
}

func mapToGetFlowsResponse(flows *getFlowsResponse) []ak.Flow {
	var flowsResp []ak.Flow

	for _, f := range flows.Results {
		flowsResp = append(flowsResp, ak.Flow{
			PK:          f.PK,
			Name:        f.Name,
			Slug:        f.Slug,
			Designation: f.Designation,
		})
	}

	return flowsResp
}

func mapToGetApplicationsByNameResponse(apps *getApplicationsResponse, name string) (a *ak.Application, err error) {
	var applicationsResp ak.Application

	var found bool
	for _, app := range apps.Results {
		if app.Name != name {
			continue
		}
		found = true

		switch app.ProviderObj.MetaModelName {
		case ProviderTypeMetaModelLDAP:
			applicationsResp.ProviderType = constants.ProviderTypeLDAP
		case ProviderTypeMetaModelOIDC:
			applicationsResp.ProviderType = constants.ProviderTypeOIDC
		case ProviderTypeMetaModelProxy:
			applicationsResp.ProviderType = constants.ProviderTypeProxy
		case ProviderTypeMetaModelRAC:
			applicationsResp.ProviderType = constants.ProviderTypeRAC
		case ProviderTypeMetaModelSAML:
			applicationsResp.ProviderType = constants.ProviderTypeSAML
		case ProviderTypeMetaModelSCIM:
			applicationsResp.ProviderType = constants.ProviderTypeSCIM
		}

		applicationsResp.PK = app.PK
		applicationsResp.Name = app.Name
		applicationsResp.Slug = app.Slug
		applicationsResp.ProviderPK = app.ProviderObj.PK
		applicationsResp.ProviderName = app.ProviderObj.Name
	}

	if !found {
		return nil, fmt.Errorf("application not found")
	}

	return &applicationsResp, nil
}
