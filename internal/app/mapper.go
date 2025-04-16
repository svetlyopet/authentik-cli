package app

import (
	"fmt"

	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/provider"
)

const (
	oidcPerProviderIssuerUrl = "%s/application/o/%s/"
	oidcGlobalIssuerUrl      = "%s/"
	oidcConfigurationUrl     = "%s/application/o/%s/.well-known/openid-configuration"
	oidcAuthorizeUrl         = "%s/application/o/authorize/"
	oidcTokenUrl             = "%s/application/o/token/"
	oidcUserinfoUrl          = "%s/application/o/userinfo/"
	oidcLogoutUrl            = "%s/application/o/%s/end-session/"
	oidcJwkUrl               = "%s/application/o/%s/jwks/"
)

func mapToGetAppWithOidcProvider(app *ak.Application, prov ak.OidcProvider) *App {
	var oidcClientSecret string

	if prov.ClientType == provider.OidcClientTypeConfidential {
		oidcClientSecret = prov.ClientSecret
	}

	authentikHost := ak.Repo.GetAuthentikTargetUrl()

	var oidcIssuer string
	if prov.IssuerMode == provider.OidcIssuerModeGlobal {
		oidcIssuer = fmt.Sprintf(oidcGlobalIssuerUrl, authentikHost)
	} else {
		oidcIssuer = fmt.Sprintf(oidcPerProviderIssuerUrl, authentikHost, app.Slug)
	}

	return &App{
		Name:         app.Name,
		ProviderType: app.ProviderType,
		OidcProvider: OidcProvider{
			ClientType:            prov.ClientType,
			ClientId:              prov.ClientId,
			ClientSecret:          oidcClientSecret,
			RedirectUris:          mapToRedirectUris(prov.RedirectUris),
			Issuer:                oidcIssuer,
			ConfigurationEndpoint: fmt.Sprintf(oidcConfigurationUrl, authentikHost, app.Slug),
			AuthorizeEndpoint:     fmt.Sprintf(oidcAuthorizeUrl, authentikHost),
			TokenEndpoint:         fmt.Sprintf(oidcTokenUrl, authentikHost),
			UserinfoEndpoint:      fmt.Sprintf(oidcUserinfoUrl, authentikHost),
			LogoutEndpoint:        fmt.Sprintf(oidcLogoutUrl, authentikHost, app.Slug),
			JwkEndpoint:           fmt.Sprintf(oidcJwkUrl, authentikHost, app.Slug),
		},
	}
}

func mapToRedirectUris(uris []ak.OidcRedirectUri) (oidcRedirectUris []OidcRedirectUri) {
	for _, uri := range uris {
		uriDetails := OidcRedirectUri{
			MatchingMode: uri.MatchingMode,
			Url:          uri.Url,
		}

		oidcRedirectUris = append(oidcRedirectUris, uriDetails)
	}

	return oidcRedirectUris
}
