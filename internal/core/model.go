package core

type Group struct {
	Name   string   `json:"name" yaml:"name"`
	Tenant string   `json:"tenant" yaml:"tenant"`
	Roles  []string `json:"roles" yaml:"roles"`
}

type User struct {
	Name        string      `json:"name" yaml:"name"`
	Email       string      `json:"email" yaml:"email"`
	IsActive    bool        `json:"is_active" yaml:"is_active"`
	IsSuperuser bool        `json:"is_superuser" yaml:"is_superuser"`
	Attributes  interface{} `json:"attributes" yaml:"attributes"`
	Groups      []string    `json:"groups" yaml:"groups"`
}
