package ak

type Role struct {
	PK   string
	Name string
}

type Group struct {
	PK   string
	Name string
}

type GroupAttributes struct {
	Tenant string
}

type User struct {
	PK         int
	Username   string
	Name       string
	Email      string
	Path       string
	IsActive   bool
	Attributes UserAttributes
}

type UserAttributes struct {
	UserType string
	Tenant   string
}
