package ak

var Repo AuthentikRepository

type AuthentikRepository interface {
	CreateRole(name string) (*Role, error)
	GetRoleByName(name string) (*Role, error)
	AssignTenantAdminPermissionsToRole(rolePK string) error
	CreateGroup(name string, roles []string, attributes map[string]string) (*Group, error)
	GetGroupByName(name string) (*Group, error)
}
