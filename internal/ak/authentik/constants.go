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
	ProviderTypeMetaModelLDAP  = "authentik_providers_ldap.ldapprovider"
	ProviderTypeMetaModelOIDC  = "authentik_providers_oauth2.oauth2provider"
	ProviderTypeMetaModelProxy = "authentik_providers_proxy.proxyprovider"
	ProviderTypeMetaModelRAC   = "authentik_providers_rac.racprovider"
	ProviderTypeMetaModelSAML  = "authentik_providers_saml.samlprovider"
	ProviderTypeMetaModelSCIM  = "authentik_providers_scim.scimprovider"
)
