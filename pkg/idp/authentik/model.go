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

type createRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type createRoleResponse struct {
	PK   string `json:"pk" binding:"required"`
	Name string `json:"name" binding:"required"`
}
