package ak

type Role struct {
	PK   string
	Name string
}

type Group struct {
	PK              string
	Name            string
	GroupAttributes GroupAttributes
	Roles           []Role
}

type GroupAttributes struct {
	Tenant string
}

type User struct {
	PK          int
	Username    string
	Name        string
	Email       string
	Path        string
	IsActive    bool
	IsSuperuser bool
	Attributes  UserAttributes
	Groups      []Group
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

// UserPatch lists the user fields to change in a partial update. A nil field
// is not sent, so authentik keeps its current value.
type UserPatch struct {
	Username *string
	Name     *string
	Email    *string
	IsActive *bool
}

// IsEmpty reports whether the patch changes nothing.
func (p UserPatch) IsEmpty() bool {
	return p.Username == nil && p.Name == nil && p.Email == nil && p.IsActive == nil
}

// GroupPatch lists the group fields to change in a partial update. A nil
// field is not sent, so authentik keeps its current value.
type GroupPatch struct {
	Name *string
}

// IsEmpty reports whether the patch changes nothing.
func (p GroupPatch) IsEmpty() bool {
	return p.Name == nil
}

// ApplicationPatch lists the application fields to change in a partial
// update. A nil field is not sent, so authentik keeps its current value.
type ApplicationPatch struct {
	Name *string
	Slug *string
}

// IsEmpty reports whether the patch changes nothing.
func (p ApplicationPatch) IsEmpty() bool {
	return p.Name == nil && p.Slug == nil
}

// OidcProviderPatch lists the OIDC provider fields to change in a partial
// update. A nil field is not sent, so authentik keeps its current value.
type OidcProviderPatch struct {
	Name              *string
	ClientType        *string
	AuthorizationFlow *string
	// RedirectUris replaces the provider's whole list of redirect URIs.
	RedirectUris *[]OidcRedirectUri
}

// IsEmpty reports whether the patch changes nothing.
func (p OidcProviderPatch) IsEmpty() bool {
	return p.Name == nil && p.ClientType == nil && p.AuthorizationFlow == nil && p.RedirectUris == nil
}
