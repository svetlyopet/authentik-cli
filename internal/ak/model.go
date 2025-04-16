package ak

type Role struct {
	PK   string
	Name string
}

type Group struct {
	PK              string
	Name            string
	GroupAttributes GroupAttributes
}

type GroupAttributes struct {
	Tenant string
}

type User struct {
	PK         int
	Username   string
	Name       string
	Email      string
	Path       string
	IsActive   bool
	Attributes UserAttributes
}

type UserAttributes struct {
	UserType string
	Tenant   string
}

type Provider struct {
	PK                 int
	Name               string
	AuthenticationFlow string
	AuthorizationFlow  string
	InvalidationFlow   string
}

type OidcProvider struct {
	Provider
	PropertyMappings     []string
	ClientType           string
	ClientId             string
	ClientSecret         string
	AccessCodeValidity   string
	AccessTokenValidity  string
	RefreshTokenValidity string
	SigningKey           string
	EncryptionKey        string
	RedirectUris         []OidcRedirectUri
	SubMode              string
	IssuerMode           string
}

type OidcRedirectUri struct {
	MatchingMode string
	Url          string
}

type Application struct {
	PK           string
	Name         string
	Slug         string
	ProviderPK   int
	ProviderType string
	ProviderName string
}

type Flow struct {
	PK          string
	Name        string
	Slug        string
	Designation string
}
