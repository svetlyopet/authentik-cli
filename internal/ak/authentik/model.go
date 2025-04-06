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
	PK          string          `json:"pk" binding:"required"`
	NumPK       int             `json:"num_pk" binding:"required"`
	Name        string          `json:"name" binding:"required"`
	IsSuperuser bool            `json:"is_superuser"`
	Parent      string          `json:"parent"`
	ParentName  string          `json:"parent_name"`
	Users       []int           `json:"users"`
	UsersObj    []userObj       `json:"users_obj" binding:"required"`
	Attributes  groupAttributes `json:"attributes"`
	Roles       []string        `json:"roles"`
	RolesObj    []roleObj       `json:"roles_obj" binding:"required"`
}

type groupAttributes struct {
	Tenant string `json:"tenant"`
}

type userObj struct {
	PK         int            `json:"pk" binding:"required"`
	Username   string         `json:"username" binding:"required"`
	Name       string         `json:"name" binding:"required"`
	IsActive   bool           `json:"is_active"`
	LastLogin  string         `json:"last_login"`
	Email      string         `json:"email"`
	Path       string         `json:"path"`
	Attributes userAttributes `json:"attributes"`
	Uid        string         `json:"uid" binding:"required"`
}

type userAttributes struct {
	UserType string `json:"userType"`
	Tenant   string `json:"tenant"`
}

type roleObj struct {
	PK   string `json:"pk" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type createRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type createOrUpdateRoleResponse struct {
	roleObj
}

type getRolesResponse struct {
	Pagination pagination `json:"pagination"`
	Results    []roleObj  `json:"results" binding:"required"`
}

type createGroupRequest struct {
	Name        string          `json:"name" binding:"required"`
	IsSuperuser bool            `json:"is_superuser"`
	Parent      string          `json:"parent"`
	Users       []int           `json:"users"`
	Attributes  groupAttributes `json:"attributes"`
	Roles       []string        `json:"roles"`
}

type createOrUpdateGroupResponse struct {
	groupsObj
}

type getGroupsResponse struct {
	Pagination pagination  `json:"pagination"`
	Results    []groupsObj `json:"results" binding:"required"`
}

type assignPermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"`
}

type createUserRequest struct {
	Username   string         `json:"username" binding:"required"`
	Name       string         `json:"name" binding:"required"`
	Email      string         `json:"email" binding:"required"`
	Path       string         `json:"path" binding:"required"`
	IsActive   bool           `json:"is_active" binding:"required"`
	Attributes userAttributes `json:"attributes"`
}

type createOrUpdateUserResponse struct {
	userObj
}

type groupUserAddRequest struct {
	PK string `json:"pk" binding:"required"`
}

type getUserResponse struct {
	Pagination pagination `json:"pagination" binding:"required"`
	Results    []userObj  `json:"results" binding:"required"`
}
