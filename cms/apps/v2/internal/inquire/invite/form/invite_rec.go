package form

import (
	"fmt"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/inquire/invite/dao"
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

type InviteRec struct {
	Date           string `json:"date"`
	UserID         int64  `json:"user_id"`
	UserName       string `json:"user_name"`
	ParentUserID   int64  `json:"parent_user_id"`
	ParentUserName string `json:"parent_user_name"`
}

func Format(dateRange [2]time.Time, _summary map[string]any, data []*dao.InviteRec) (summary map[string]any, result []*InviteRec) {
	summary = _summary
	if summary != nil {
		summary["amount"] = util.ConvertAmount2Decimal(summary["amount"])
	}

	for _, item := range data {
		result = append(result, &InviteRec{
			Date:           item.Date,
			UserID:         item.UserID,
			UserName:       item.UserName,
			ParentUserID:   item.ParentUserID,
			ParentUserName: item.ParentUserName,
		})
	}

	return
}

func Format2Excel(dateRange [2]time.Time, _data []*dao.InviteRec) (excelModel *excel.Excel[*InviteRec], err error) {
	_, data := Format(dateRange, nil, _data)

	reflectMap := map[string]func(source *InviteRec) any{
		"日期":     func(source *InviteRec) any { return source.Date },
		"被邀用户ID": func(source *InviteRec) any { return source.UserID },
		"被邀用户昵称": func(source *InviteRec) any { return source.UserName },
		"邀请用户ID": func(source *InviteRec) any { return source.ParentUserID },
		"邀请用户昵称": func(source *InviteRec) any { return source.ParentUserName },
	}

	excelModel = &excel.Excel[*InviteRec]{
		FileName:   fmt.Sprintf("invite_rec_%s-%s", dateRange[0].Format(pkg.FILE_DATE_FORMAT), dateRange[1].Format(pkg.FILE_DATE_FORMAT)),
		SheetNames: []string{"邀请记录"},
		SheetNameWithHead: map[string][]string{
			"邀请记录": {
				"日期",
				"被邀用户ID", "被邀用户昵称", "邀请用户ID", "邀请用户昵称",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*InviteRec]{
			"邀请记录": data,
		},
		ReflectMap: map[string]excel.RowReflect[*InviteRec]{
			"邀请记录": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
