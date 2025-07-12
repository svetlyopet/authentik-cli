package authentik

type pagination struct {
	Next       int `json:"next"`
	Previous   int `json:"previous"`
	Count      int `json:"count"`
	Current    int `json:"current"`
	TotalPages int `json:"total_pages"`
	StartIndex int `json:"start_index"`
	EndIndex   int `json:"end_index"`
}

type groupObj struct {
	PK          string          `json:"pk" binding:"required"`
	NumPK       int             `json:"num_pk" binding:"required"`
	Name        string          `json:"name" binding:"required"`
	IsSuperuser bool            `json:"is_superuser"`
	Parent      string          `json:"parent,omitempty"`
	ParentName  string          `json:"parent_name,omitempty"`
	Users       []int           `json:"users"`
	UsersObj    []userObj       `json:"users_obj" binding:"required"`
	Attributes  groupAttributes `json:"attributes"`
	Roles       []string        `json:"roles"`
	RolesObj    []roleObj       `json:"roles_obj" binding:"required"`
}

type groupAttributes struct {
	Tenant string `json:"tenant,omitempty"`
}

type userObj struct {
	PK          int            `json:"pk" binding:"required"`
	Username    string         `json:"username" binding:"required"`
	Name        string         `json:"name" binding:"required"`
	IsActive    bool           `json:"is_active" binding:"required"`
	IsSuperuser bool           `json:"is_superuser" binding:"required"`
	LastLogin   string         `json:"last_login"`
	Email       string         `json:"email" binding:"required"`
	Path        string         `json:"path" binding:"required"`
	Attributes  userAttributes `json:"attributes"`
	Uid         string         `json:"uid" binding:"required"`
	Groups      []string       `json:"groups"`
	GroupsObj   []groupObj     `json:"groups_obj" binding:"required"`
}

type userAttributes struct {
	UserType string `json:"userType,omitempty"`
	Tenant   string `json:"tenant,omitempty"`
}

type roleObj struct {
	PK   string `json:"pk" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// baseProvider contains all the components of a provider which are shared between all provider types
type baseProvider struct {
	PK                                 int      `json:"pk"  binding:"required"`
	Name                               string   `json:"name"  binding:"required"`
	AuthenticationFlow                 string   `json:"authentication_flow"  binding:"required"`
	AuthorizationFlow                  string   `json:"authorization_flow"  binding:"required"`
	InvalidationFlow                   string   `json:"invalidation_flow"  binding:"required"`
	PropertyMappings                   []string `json:"property_mappings"  binding:"required"`
	Component                          string   `json:"component"  binding:"required"`
	AssignedApplicationSlug            string   `json:"assigned_application_slug"  binding:"required"`
	AssignedApplicationName            string   `json:"assigned_application_name"  binding:"required"`
	AssignedBackchannelApplicationSlug string   `json:"assigned_backchannel_application_slug"`
	AssignedBackchannelApplicationName string   `json:"assigned_backchannel_application_name"`
	VerboseName                        string   `json:"verbose_name"  binding:"required"`
	VerboseNamePlural                  string   `json:"verbose_name_plural"  binding:"required"`
	MetaModelName                      string   `json:"meta_model_name"  binding:"required"`
}

type oidcProvider struct {
	baseProvider
	ClientType             string            `json:"client_type"  binding:"required"`
	ClientId               string            `json:"client_id"  binding:"required"`
	ClientSecret           string            `json:"client_secret"`
	AccessCodeValidity     string            `json:"access_code_validity"  binding:"required"`
	AccessTokenValidity    string            `json:"access_token_validity"  binding:"required"`
	RefreshTokenValidity   string            `json:"refresh_token_validity"  binding:"required"`
	IncludeClaimsInIdToken bool              `json:"include_claims_in_id_token"  binding:"required"`
	SigningKey             string            `json:"signing_key"  binding:"required"`
	EncryptionKey          string            `json:"encryption_key"`
	RedirectUris           []oidcRedirectUri `json:"redirect_uris"  binding:"required"`
	SubMode                string            `json:"sub_mode"  binding:"required"`
	IssuerMode             string            `json:"issuer_mode"  binding:"required"`
	JWKSSources            []string          `json:"jwk_sources"`
}

type application struct {
	PK                      string         `json:"pk" binding:"required"`
	Name                    string         `json:"name" binding:"required"`
	Slug                    string         `json:"slug" binding:"required"`
	Provider                int            `json:"provider"`
	ProviderObj             baseProvider   `json:"provider_obj" binding:"required"`
	BackchannelProviders    []int          `json:"backchannel_providers"`
	BackchannelProvidersObj []baseProvider `json:"backchannel_providers_obj" binding:"required"`
	LaunchUrl               string         `json:"launch_url" binding:"required"`
	OpenInNewTab            bool           `json:"open_in_new_tab"`
	MetaLaunchUrl           string         `json:"meta_launch_url"`
	MetaIcon                string         `json:"meta_icon" binding:"required"`
	MetaDescription         string         `json:"meta_description"`
	MetaPublisher           string         `json:"meta_publisher"`
	PolicyEngineMode        string         `json:"policy_engine_mode"`
	Group                   string         `json:"group"`
}

type flow struct {
	PK                      string   `json:"pk" binding:"required"`
	PolicyBindingModelPtrId string   `json:"policyindingmodel_ptr_id" binding:"required"`
	Name                    string   `json:"name" binding:"required"`
	Slug                    string   `json:"slug" binding:"required"`
	Title                   string   `json:"title" binding:"required"`
	Designation             string   `json:"designation" binding:"required"`
	Background              string   `json:"background" binding:"required"`
	Stages                  []string `json:"stages" binding:"required"`
	Policies                []string `json:"policies" binding:"required"`
	CacheCount              int      `json:"cache_count" binding:"required"`
	PolicyEngineMode        string   `json:"policy_engine_mode"`
	CompatabilityMode       bool     `json:"compatability_mode"`
	ExportUrl               string   `json:"export_url" binding:"required"`
	Layout                  string   `json:"layout"`
	DeniedAction            string   `json:"denied_action"`
	Authentication          string   `json:"authentication"`
}

type createRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type getRoleResponse struct {
	roleObj
}

type getRolesResponse struct {
	Pagination pagination `json:"pagination"`
	Results    []roleObj  `json:"results" binding:"required"`
}

type createGroupRequest struct {
	Name        string          `json:"name" binding:"required"`
	IsSuperuser bool            `json:"is_superuser"`
	Parent      string          `json:"parent,omitempty"`
	Users       []int           `json:"users,omitempty"`
	Attributes  groupAttributes `json:"attributes"`
	Roles       []string        `json:"roles,omitempty"`
}

type getGroupResponse struct {
	groupObj
}

type getGroupsResponse struct {
	Pagination pagination `json:"pagination"`
	Results    []groupObj `json:"results" binding:"required"`
}

type assignPermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"`
}

type createUserRequest struct {
	Username   string         `json:"username" binding:"required"`
	Name       string         `json:"name" binding:"required"`
	Email      string         `json:"email" binding:"required"`
	Path       string         `json:"path" binding:"required"`
	IsActive   bool           `json:"is_active" binding:"required"`
	Attributes userAttributes `json:"attributes"`
}

type getUserResponse struct {
	userObj
}

type groupUserAddRequest struct {
	PK string `json:"pk" binding:"required"`
}

type getUsersResponse struct {
	Pagination pagination `json:"pagination" binding:"required"`
	Results    []userObj  `json:"results" binding:"required"`
}

type createOrUpdateOidcProviderRequest struct {
	Name                   string            `json:"name" binding:"required"`
	AuthenticationFlow     string            `json:"authentication_flow" binding:"required"`
	AuthorizationFlow      string            `json:"authorization_flow" binding:"required"`
	InvalidationFlow       string            `json:"invalidation_flow" binding:"required"`
	PropertyMappings       []string          `json:"property_mappings,omitempty"`
	ClientType             string            `json:"client_type" binding:"required"`
	AccessCodeValidity     string            `json:"access_code_validity,omitempty"`
	AccessTokenValidity    string            `json:"access_token_validity,omitempty"`
	RefreshTokenValidity   string            `json:"refresh_token_validity,omitempty"`
	IncludeClaimsInIdToken bool              `json:"include_claims_in_id_token,omitempty"`
	SigningKey             string            `json:"signing_key,omitempty"`
	EncryptionKey          string            `json:"encryption_key,omitempty"`
	RedirectUris           []oidcRedirectUri `json:"redirect_uris" binding:"required"`
	SubMode                string            `json:"sub_mode" binding:"required"`
	IssuerMode             string            `json:"issuer_mode" binding:"required"`
	JwtFederationSources   []string          `json:"jwt_federation_sources,omitempty"`
	JwtFederationProviders []int             `json:"jwt_federation_providers,omitempty"`
}

type oidcRedirectUri struct {
	MatchingMode string `json:"matching_mode" binding:"required"`
	Url          string `json:"url" binding:"required"`
}

type getOidcProviderResponse struct {
	oidcProvider
}

type createOrUpdateApplicationRequest struct {
	Name                 string `json:"name" binding:"required"`
	Slug                 string `json:"slug" binding:"required"`
	Provider             int    `json:"provider"`
	BackchannelProviders []int  `json:"backchannel_providers,omitempty"`
	OpenInNewTab         bool   `json:"open_in_new_tab,omitempty"`
	MetaLaunchUrl        string `json:"meta_launch_url,omitempty"`
	MetaDescription      string `json:"meta_description,omitempty"`
	PolicyEngineMode     string `json:"policy_engine_mode,omitempty"`
	Group                string `json:"group,omitempty"`
}

type getApplicationResponse struct {
	application
}

type getFlowsResponse struct {
	Pagination pagination `json:"pagination" binding:"required"`
	Results    []flow     `json:"results" binding:"required"`
}

type getApplicationsResponse struct {
	Pagination pagination    `json:"pagination" binding:"required"`
	Results    []application `json:"results" binding:"required"`
}
