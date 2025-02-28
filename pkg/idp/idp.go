package idp

type AuthentikRepository interface {
	CreateRole(name string) (*Role, error)
	CreateGroup(name string, roles []string, attributes map[string]string) (*Group, error)
}
