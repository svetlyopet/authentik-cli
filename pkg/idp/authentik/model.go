package authentik

type pagination struct {
	Next       int `json:"next"`
	Previous   int `json:"previous"`
	Count      int `json:"count"`
	Current    int `json:"current"`
	TotalPages int `json:"total_pages"`
	StartIndex int `json:"start_index"`
	EndIndex   int `json:"end_intex"`
}

type groupsObj struct {
	PK          string      `json:"pk" binding:"required"`
	NumPK       int         `json:"num_pk" binding:"required"`
	Name        string      `json:"name" binding:"required"`
	IsSuperuser bool        `json:"is_superuser"`
	Parent      string      `json:"parent"`
	ParentName  string      `json:"parent_name"`
	Users       []int       `json:"users"`
	UsersObj    []usersObj  `json:"users_obj" binding:"required"`
	Attributes  interface{} `json:"attributes"`
	Roles       []string    `json:"roles"`
	RolesObj    []rolesObj  `json:"roles_obj" binding:"required"`
}

type usersObj struct {
	PK         int         `json:"pk" binding:"required"`
	Username   string      `json:"username" binding:"required"`
	Name       string      `json:"name" binding:"required"`
	IsActive   bool        `json:"is_active"`
	LastLogin  string      `json:"last_login"`
	Email      string      `json:"email"`
	Attributes interface{} `json:"attributes"`
	Uid        string      `json:"uid" binding:"required"`
}

type rolesObj struct {
	PK   string `json:"pk" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type createRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type createRoleResponse struct {
	rolesObj
}

type getRolesResponse struct {
	Pagination pagination `json:"pagination"`
	Results    []rolesObj `json:"results" binding:"required"`
}

type createGroupRequest struct {
	Name        string      `json:"name" binding:"required"`
	IsSuperuser bool        `json:"is_superuser"`
	Parent      string      `json:"parent"`
	Users       []int       `json:"users"`
	Attributes  interface{} `json:"attributes"`
	Roles       []string    `json:"roles"`
}

type createGroupResponse struct {
	groupsObj
}

type getGroupsResponse struct {
	Pagination pagination  `json:"pagination"`
	Results    []groupsObj `json:"results" binding:"required"`
}

type assignPermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"`
}
