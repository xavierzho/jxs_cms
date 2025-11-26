package form

import (
	"fmt"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/report/recall/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type GenerateRequest struct {
	iForm.DateRangeRequest
}

type UserType string

const (
	UserType_None       UserType = ""
	UserType_User       UserType = "user"
	UserType_ParentUser UserType = "parent_user"
)

type ListRequest struct {
	app.Pager
	AllRequest
}

func (q *ListRequest) Parse() (dateRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	q.Pager.Parse()

	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateRangeRequest
	cForm.UserInfoRequest
	UserType UserType `form:"user_type"`
}

func (q *AllRequest) Parse() (dateRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	if err = q.Valid(); err != nil {
		return dateRange, nil, err
	}

	if dateRange, err = q.DateRangeRequest.Parse(); err != nil {
		return
	}

	if queryParams, err = q.UserInfoRequest.Parse(); err != nil {
		return
	}

	return dateRange, queryParams, nil
}

func (q *AllRequest) Valid() error {
	switch q.UserType {
	case UserType_None:
	case UserType_User, UserType_ParentUser:
	default:
		return fmt.Errorf("not expected UserType: %v", q.UserType)
	}

	return nil
}

type Recall struct {
	Date           string          `json:"date"`
	UserID         int64           `json:"user_id"`
	UserName       string          `json:"user_name"`
	ParentUserID   int64           `json:"parent_user_id"`
	ParentUserName string          `json:"parent_user_name"`
	Amount         decimal.Decimal `json:"amount"`
	Point          string          `json:"point"`
}

func Format(dateRange [2]time.Time, _summary map[string]any, data []*dao.Recall) (summary map[string]any, result []*Recall) {
	summary = _summary
	if summary != nil {
		summary["amount"] = util.ConvertAmount2Decimal(summary["amount"])
		summary["point"] = convert.GetDecimal(summary["point"]).Div(decimal.NewFromInt(10))
	}

	for _, item := range data {
		result = append(result, &Recall{
			Date:           item.Date,
			UserID:         item.UserID,
			UserName:       item.UserName,
			ParentUserID:   item.ParentUserID,
			ParentUserName: item.ParentUserName,
			Amount:         util.ConvertAmount2Decimal(item.Amount),
			Point:          convert.GetString(item.Point / 10),
		})
	}

	return
}

func Format2Excel(dateRange [2]time.Time, _data []*dao.Recall) (excelModel *excel.Excel[*Recall], err error) {
	_, data := Format(dateRange, nil, _data)
	reflectMap := map[string]func(source *Recall) any{
		"日期":        func(source *Recall) any { return source.Date },
		"被召回用户ID":   func(source *Recall) any { return source.UserID },
		"被召回用户昵称":   func(source *Recall) any { return source.UserName },
		"召回用户ID":    func(source *Recall) any { return source.ParentUserID },
		"召回用户昵称":    func(source *Recall) any { return source.ParentUserName },
		"被召回用户抽赏金额": func(source *Recall) any { return source.Amount },
		"召回用户获得欧气值": func(source *Recall) any { return source.Point },
	}
	excelModel = &excel.Excel[*Recall]{
		FileName:   fmt.Sprintf("recall_%s-%s", dateRange[0].Format(pkg.FILE_DATE_FORMAT), dateRange[1].Format(pkg.FILE_DATE_FORMAT)),
		SheetNames: []string{"召回报表"},
		SheetNameWithHead: map[string][]string{
			"召回报表": {
				"日期",
				"被召回用户ID", "被召回用户昵称", "召回用户ID", "召回用户昵称",
				"被召回用户抽赏金额", "召回用户获得欧气值",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*Recall]{
			"召回报表": data,
		},
		ReflectMap: map[string]excel.RowReflect[*Recall]{
			"召回报表": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
