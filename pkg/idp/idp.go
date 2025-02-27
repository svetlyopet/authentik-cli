package idp

type AuthentikRepository interface {
	CreateRole(name string) (*Role, error)
}
