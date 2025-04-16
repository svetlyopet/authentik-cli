package app

type App struct {
	Name         string       `json:"name" yaml:"name"`
	ProviderType string       `json:"provider_type" yaml:"provider_type"`
	OidcProvider OidcProvider `json:"oidc_provider,omitempty" yaml:"oidc_provider,omitempty"`
}

type OidcProvider struct {
	ClientType            string            `json:"client_type" yaml:"client_type"`
	ClientId              string            `json:"client_id" yaml:"client_id"`
	ClientSecret          string            `json:"client_secret,omitempty" yaml:"client_secret,omitempty"`
	RedirectUris          []OidcRedirectUri `json:"redirect_uris" yaml:"redirect_uris"`
	Issuer                string            `json:"issuer" yaml:"issuer"`
	ConfigurationEndpoint string            `json:"configuration_endpoint" yaml:"configuration_endpoint"`
	AuthorizeEndpoint     string            `json:"authorize_endpoint" yaml:"authorize_endpoint"`
	TokenEndpoint         string            `json:"token_endpoint" yaml:"token_endpoint"`
	UserinfoEndpoint      string            `json:"userinfo_endpoint" yaml:"userinfo_endpoint"`
	LogoutEndpoint        string            `json:"logout_endpoint" yaml:"logout_endpoint"`
	JwkEndpoint           string            `json:"jwk_endpoint" yaml:"jwk_endpoint"`
}

type OidcRedirectUri struct {
	MatchingMode string `json:"matching_mode" yaml:"matching_mode"`
	Url          string `json:"url" yaml:"url"`
}
