package provider

const (
	OidcClientTypeConfidential = "confidential"
	OidcClientTypePublic       = "public"
)

const (
	OidcRedirectUriMatchingModeStrict = "strict"
	OidcRedirectUriMatchingModeRegex  = "regex"
)

const (
	OidcSubModeDefault      = OidcSubModeUserEmail
	OidcSubModeHashedUserId = "hashed_user_id"
	OidcSubModeUserId       = "user_id"
	OidcSubModeUserUuid     = "user_uuid"
	OidcSubModeUserUsername = "user_username"
	OidcSubModeUserEmail    = "user_email"
	OidcSubModeUserUpn      = "user_upn"
)

const (
	OidcIssuerModeDefault     = OidcIssuerModePerProvider
	OidcIssuerModeGlobal      = "global"
	OidcIssuerModePerProvider = "per_provider"
)

const (
	OidcAccessCodeValidityDefault   = "minutes=5"
	OidcAccessTokenValidityDefault  = "minutes=30"
	OidcRefreshTokenValidityDefault = "hours=24"
)

const (
	OidcConsentTypeExplicit = "explicit"
	OidcConsentTypeImplicit = "implicit"
)

// Slugs of the flows created by authentik's default blueprints.
const (
	FlowSlugAuthentication               = "default-authentication-flow"
	FlowSlugAuthorizationExplicitConsent = "default-provider-authorization-explicit-consent"
	FlowSlugAuthorizationImplicitConsent = "default-provider-authorization-implicit-consent"
	FlowSlugInvalidation                 = "default-invalidation-flow"
)
