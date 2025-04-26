package core

type Group struct {
	Name   string   `json:"name" yaml:"name"`
	Tenant string   `json:"tenant,omitempty" yaml:"tenant,omitempty"`
	Roles  []string `json:"roles,omitempty" yaml:"roles,omitempty"`
}
