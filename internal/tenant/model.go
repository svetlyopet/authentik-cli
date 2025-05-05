package tenant

type Tenant struct {
	Name  string `json:"name" yaml:"name"`
	Group string `json:"group" yaml:"group"`
	Role  string `json:"role" yaml:"role"`
}
