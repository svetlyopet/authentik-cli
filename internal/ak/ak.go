package ak

var Repo AuthentikRepository

type AuthentikRepository interface {
	GetAuthentikTargetUrl() string

	AuthentikCore
	AuthentikRbac
	AuthentikProvider
}

type AuthentikCore interface {
	CreateUser(usr User) (*User, error)
	GetUserByUsername(username string) (*User, error)
	DeleteUser(userPK string) error
	CreateGroup(name string, roles []string, attributes GroupAttributes) (*Group, error)
	GetGroup(uuid string) (*Group, error)
	GetGroupByName(name string) (*Group, error)
	DeleteGroup(uuid string) error
	AddUserToGroup(userPK int, uuid string) error
	UpdateUser(userPK int, patch UserPatch) (*User, error)
	UpdateGroup(uuid string, patch GroupPatch) (*Group, error)
}

type AuthentikRbac interface {
	CreateRole(name string) (*Role, error)
	GetRoleByName(name string) (*Role, error)
	DeleteRole(uuid string) error
	AssignViewPermissionsToTenantRole(rolePK string) error
}

type AuthentikProvider interface {
	CreateOidcProvider(provider OidcProvider) (*OidcProvider, error)
	CreateApplication(name, slug string, providerPK int) (*Application, error)
	GetFlowBySlug(slug string) (*Flow, error)
	GetApplicationByName(name string) (*Application, error)
	DeleteApplication(slug string) error
	DeleteProvider(pk int) error
	GetOidcProvider(pk int) (*OidcProvider, error)
	UpdateApplication(slug string, patch ApplicationPatch) (*Application, error)
	UpdateOidcProvider(pk int, patch OidcProviderPatch) (*OidcProvider, error)
}
