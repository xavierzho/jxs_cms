package form

import "data_backend/pkg/database"

type RoleCreateRequest struct {
	Name             string   `form:"name" binding:"required,min=2,max=100"`
	PermissionIDList []uint32 `form:"permission_id_list[]" binding:"required"`
}

type RoleListRequest struct {
	Name string `form:"name" binding:"max=100"`
}

func (q *RoleListRequest) Parse() (queryParams database.QueryWhereGroup) {
	queryParams = make([]database.QueryWhere, 0, 1)
	if q.Name != "" {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "name LIKE ?",
			Value:  []interface{}{"%" + q.Name + "%"},
		})
	}

	return queryParams
}

type RoleUpdateRequest struct {
	Name             string   `form:"name" binding:"required,min=2,max=100"`
	PermissionIDList []uint32 `form:"permission_id_list[]" binding:"required"`
}
