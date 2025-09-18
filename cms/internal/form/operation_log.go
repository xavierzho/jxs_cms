package form

import (
	"errors"

	"data_backend/pkg/database"
)

type OperationLogListRequest struct {
	UserID   uint32 `json:"user_id" form:"user_id" `
	ModuleID string `json:"module_id" form:"module_id"`
	Module   string `json:"module" form:"module"`
}

func (r *OperationLogListRequest) Parse() ([]database.QueryWhere, error) {
	queryParams := make([]database.QueryWhere, 0, 3)
	if r.UserID != 0 {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "user_id = ?",
			Value:  []interface{}{r.UserID},
		})
	}
	if r.Module != "" {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "module_name = ?",
			Value:  []interface{}{r.Module},
		})
	}
	if r.ModuleID != "" {
		if r.Module == "" {
			return nil, errors.New("module is required")
		}
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "module_id = ?",
			Value:  []interface{}{r.ModuleID},
		})
	}

	return queryParams, nil
}
