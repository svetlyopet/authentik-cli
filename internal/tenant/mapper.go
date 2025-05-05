package tenant

func MapToGetTenant(name, group, role string) *Tenant {
	return &Tenant{
		Name:  name,
		Group: group,
		Role:  role,
	}
}
