package constants

// CfgFilename represents the default config filename
const CfgFilename = ".authentik-cli"

// TenantAdminRbacRoleNamePattern represents the naming pattern for naming the RBAC role for a tenant
const TenantAdminRbacRoleNamePattern = "%s-tenant-admin"

// TenantAdminGroupNamePattern represents the naming pattern for naming the admin group for a tenant
const TenantAdminGroupNamePattern = "%s-tenant-admin"

// ObjectTypeRole holds the name Role object
const ObjectTypeRole = "role"

// ObjectTypeGroup holds the name of the Group object
const ObjectTypeGroup = "group"

// ObjectTypeUser holds the name of the User object
const ObjectTypeUser = "user"

const ActionCreated = "created"
const ActionChanged = "changed"
const ActionUnchanged = "unchanged"
const ActionDeleted = "deleted"
