package form

import (
	"regexp"
	"strconv"
	"strings"

	"data_backend/pkg/database"
	"data_backend/pkg/encrypt/md5"
)

type UserInfoRequest struct {
	UserID   int64  `form:"user_id"`
	UserIDS  string `form:"user_ids"`
	UserName string `form:"user_name"`
	Tel      string `form:"tel"`
	IsAdmin  *bool  `form:"is_admin"`
	Channel  int32  `form:"channel"`
}

func (q *UserInfoRequest) Parse() (queryParams database.QueryWhereGroup, err error) {
	if q.UserIDS != "" {
		// 只保留数字和英文逗号
		re := regexp.MustCompile(`[^0-9,]+`)
		q.UserIDS = re.ReplaceAllString(q.UserIDS, "")
		// 分隔并转换成int
		ids := strings.Split(q.UserIDS, ",")
		var idsIntArr []int64
		for _, idStr := range ids {
			id, _ := strconv.ParseInt(idStr, 10, 64)
			idsIntArr = append(idsIntArr, id)
		}
		if len(ids) == 1 {
			queryParams = append(queryParams, database.QueryWhere{
				Prefix: "u.id = ?",
				Value:  []any{idsIntArr[0]},
			})
		} else {
			queryParams = append(queryParams, database.QueryWhere{
				Prefix: "u.id in ?",
				Value:  []any{idsIntArr},
			})
		}
	}

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
			Prefix: "u.role = ?",
			Value:  []any{isAdmin},
		})
	}

	if q.Channel != 0 {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "u.channel_id = ?",
			Value:  []any{q.Channel},
		})
	}

	return
}
