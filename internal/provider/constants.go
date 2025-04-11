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
