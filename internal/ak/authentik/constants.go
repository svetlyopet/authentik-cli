package authentik

const (
	AccessAdminInterfacePerm = "authentik_rbac.access_admin_interface"
	ViewSystemInfoPerm       = "authentik_rbac.view_system_info"
	ViewBrandPerm            = "authentik_brands.view_brand"
	ViewFlowPerm             = "authentik_flows.view_flow"
	ViewOutpostPerm          = "authentik_outposts.view_outpost"
	ViewApplicationPerm      = "authentik_core.view_application"
	ViewProviderPerm         = "authentik_core.view_provider"
	ViewEventPerm            = "authentik_events.view_event"
)

const (
	OidcClientTypeConfidential = "confidential"
	OidcClientTypePublic       = "public"
)

const (
	OidcRedirectUriMatchingModeStrict = "strict"
	OidcRedirectUriMatchingModeRegex  = "regex"
)

const (
	OidcSubModeHashedUserId = "hashed_user_id"
	OidcSubModeUserId       = "user_id"
	OidcSubModeUserUuid     = "user_uuid"
	OidcSubModeUserUsername = "user_username"
	OidcSubModeUserEmail    = "user_email"
	OidcSubModeUserUpn      = "user_upn"
)

const (
	OidcIssuerModeGlobal      = "global"
	OidcIssuerModePerProvider = "per_provider"
)
