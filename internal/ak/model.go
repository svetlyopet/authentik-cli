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

type OidcProvider struct {
	PK                   int
	Name                 string
	AuthenticationFlow   string
	AuthorizationFlow    string
	InvalidationFlow     string
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
