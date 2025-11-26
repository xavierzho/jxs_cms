package form

import (
	"fmt"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/inquire/recall/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"
)

type UserType string

const (
	UserType_None       UserType = ""
	UserType_User       UserType = "user"
	UserType_ParentUser UserType = "parent_user"
)

type GenerateRequest struct {
	iForm.DateRangeRequest
}

type ListRequest struct {
	app.Pager
	AllRequest
}

func (q *ListRequest) Parse() (dateRange [2]time.Time, err error) {
	q.Pager.Parse()

	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateRangeRequest
	cForm.UserInfoRequest
	UserType UserType `form:"user_type"`
}

func (q *AllRequest) Parse() (dateRange [2]time.Time, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	if dateRange, err = q.DateRangeRequest.Parse(); err != nil {
		return
	}

	return dateRange, nil
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

type RecallSvc struct {
	Date           string `json:"date"`
	UserID         int64  `json:"user_id"`
	UserName       string `json:"user_name"`
	ParentUserID   int64  `json:"parent_user_id"`
	ParentUserName string `json:"parent_user_name"`
}

func Format(dateRange [2]time.Time, _summary map[string]any, data []*dao.RecallSvc) (summary map[string]any, result []*RecallSvc) {
	summary = _summary
	if summary != nil {
		summary["amount"] = util.ConvertAmount2Decimal(summary["amount"])
	}

	for _, item := range data {
		result = append(result, &RecallSvc{
			Date:           item.Date,
			UserID:         item.UserID,
			UserName:       item.UserName,
			ParentUserID:   item.ParentUserID,
			ParentUserName: item.ParentUserName,
		})
	}

	return
}

func Format2Excel(dateRange [2]time.Time, _data []*dao.RecallSvc) (excelModel *excel.Excel[*RecallSvc], err error) {
	_, data := Format(dateRange, nil, _data)

	reflectMap := map[string]func(source *RecallSvc) any{
		"日期":      func(source *RecallSvc) any { return source.Date },
		"被召回用户ID": func(source *RecallSvc) any { return source.UserID },
		"被召回用户昵称": func(source *RecallSvc) any { return source.UserName },
		"召回用户ID":  func(source *RecallSvc) any { return source.ParentUserID },
		"召回用户昵称":  func(source *RecallSvc) any { return source.ParentUserName },
	}

	excelModel = &excel.Excel[*RecallSvc]{
		FileName:   fmt.Sprintf("recall_%s-%s", dateRange[0].Format(pkg.FILE_DATE_FORMAT), dateRange[1].Format(pkg.FILE_DATE_FORMAT)),
		SheetNames: []string{"召回记录"},
		SheetNameWithHead: map[string][]string{
			"召回记录": {
				"日期",
				"被召回用户ID", "被召回用户昵称", "召回用户ID", "召回用户昵称",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*RecallSvc]{
			"召回记录": data,
		},
		ReflectMap: map[string]excel.RowReflect[*RecallSvc]{
			"召回记录": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
