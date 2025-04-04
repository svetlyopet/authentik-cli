package ak

var Repo AuthentikRepository

type AuthentikRepository interface {
	CreateRole(name string) (*Role, error)
	GetRoleByName(name string) (*Role, error)
	DeleteRole(uuid string) error
	AssignViewPermissionsToTenantRole(rolePK string) error
	CreateGroup(name string, roles []string, attributes GroupAttributes) (*Group, error)
	GetGroupByName(name string) (*Group, error)
	DeleteGroup(uuid string) error
	AddUserToGroup(userPK int, uuid string) error
	CreateUser(usr User) (*User, error)
	GetUserByUsername(username string) (*User, error)
	DeleteUser(userPK string) error
}
