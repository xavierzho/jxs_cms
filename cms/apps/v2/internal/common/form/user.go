package form

import (
	"data_backend/pkg/database"
	"data_backend/pkg/encrypt/md5"
)

type UserInfoRequest struct {
	UserID   int64  `form:"user_id"`
	UserName string `form:"user_name"`
	Tel      string `form:"tel"`
	IsAdmin  *bool  `form:"is_admin"`
}

func (q *UserInfoRequest) Parse() (queryParams database.QueryWhereGroup, err error) {
	if q.UserID != 0 {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "u.id = ?",
			Value:  []any{q.UserID},
		})
	}

	if q.UserName != "" {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "u.nickname = ?",
			Value:  []any{q.UserName},
		})
	}

	if q.Tel != "" {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "u.phone_num_md5 = ?",
			Value:  []any{md5.EncodeMD5(q.Tel)},
		})
	}

	if q.IsAdmin != nil {
		var isAdmin = 0
		if *q.IsAdmin {
			isAdmin = 1
		}
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "u.is_admin = ?",
			Value:  []any{isAdmin},
		})
	}

	return
}
