package form

import (
	"fmt"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
)

type Request struct {
	app.Pager
	AllRequest
}

func (q *Request) Parse() (dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	q.Pager.Parse()

	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateTimeRangeRequest
	cForm.UserInfoRequest
}

func (q *AllRequest) Parse() (dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	if dateTimeRange, err = q.DateTimeRangeRequest.Parse(); err != nil {
		return
	}

	if queryParams, err = q.UserInfoRequest.Parse(); err != nil {
		return
	}

	return
}

type SignIn struct {
	CreatedAt  string `json:"created_at"`
	UserID     int64  `json:"user_id"`
	UserName   string `json:"user_name"`
	Type       string `json:"type"`
	Value      int64  `json:"value"`
	DayNo      uint32 `json:"day_no"`
	SignInType string `json:"sign_in_type"`
}

const (
	// 奖励类型
	TypeCostAwardPoint = "CostAwardPoint"
	TypeCoupon         = "Coupon"
	TypeItem           = "Item"

	// 欧气值转换比例
	CostAwardPointRatio = 10

	UnknownType     = "未知类型"
	UnknownSignType = "未知签到类型"
)

// 格式化签到奖励类型
func getAwardTypeText(typeStr string) string {
	typeMapping := map[string]string{
		TypeCostAwardPoint: "欧气值",
		TypeCoupon:         "优惠券",
		TypeItem:           "物品",
	}
	if text, exists := typeMapping[typeStr]; exists {
		return text
	}
	return UnknownType
}

// 格式化签到类型
func getSignInTypeText(signInType int64) string {
	signInTypeMapping := map[int64]string{
		1: "登录签到",
		2: "额外签到",
	}
	if text, exists := signInTypeMapping[signInType]; exists {
		return text
	}
	return UnknownSignType
}

func Format(data []map[string]any) (result []*SignIn) {
	for _, item := range data {
		typeStr := convert.GetString(item["type"])

		// 处理奖励值
		var value int64
		if typeStr == TypeCostAwardPoint {
			value = convert.GetInt64(item["value"]) / CostAwardPointRatio
		} else {
			value = convert.GetInt64(item["value"])
		}

		result = append(result, &SignIn{
			CreatedAt:  convert.GetString(item["created_at"]),
			UserID:     convert.GetInt64(item["user_id"]),
			UserName:   convert.GetString(item["user_name"]),
			SignInType: getSignInTypeText(convert.GetInt64(item["sign_in_type"])),
			DayNo:      convert.GetUint32(item["day_no"]),
			Type:       getAwardTypeText(typeStr),
			Value:      value,
		})
	}
	return
}

func Format2Excel(dateTimeRange [2]time.Time, _data []map[string]any) (excelModel *excel.Excel[*SignIn], err error) {
	data := Format(_data)
	reflectMap := map[string]func(source *SignIn) any{
		"时间":   func(source *SignIn) any { return source.CreatedAt },
		"用户ID": func(source *SignIn) any { return source.UserID },
		"用户昵称": func(source *SignIn) any { return source.UserName },
		"签到类型": func(source *SignIn) any { return source.SignInType },
		"第几天":  func(source *SignIn) any { return source.DayNo },
		"奖励类型": func(source *SignIn) any { return source.Type },
		"奖励值":  func(source *SignIn) any { return source.Value },
	}
	excelModel = &excel.Excel[*SignIn]{
		FileName:   fmt.Sprintf("sign_in_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"签到日志"},
		SheetNameWithHead: map[string][]string{
			"签到日志": {
				"时间", "用户ID", "用户昵称",
				"签到类型", "第几天", "奖励类型", "奖励值",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*SignIn]{
			"签到日志": data,
		},
		ReflectMap: map[string]excel.RowReflect[*SignIn]{
			"签到日志": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
