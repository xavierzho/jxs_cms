package form

import (
	"fmt"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/report/order/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type GenerateRequest struct {
	iForm.DateRangeRequest
}

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

type DeliveryOrder struct {
	Date           string          `json:"date"`
	UserID         int64           `json:"user_id"`
	UserName       string          `json:"user_name"`
	ShowPrice      decimal.Decimal `json:"show_price"`
	InnerPrice     decimal.Decimal `json:"inner_price"`
	RecyclingPrice decimal.Decimal `json:"recycling_price"`
}

func Format(dateRange [2]time.Time, _summary map[string]any, data []*dao.DeliveryOrder) (summary map[string]any, result []*DeliveryOrder) {
	summary = _summary
	if summary != nil {
		summary["show_price"] = util.ConvertAmount2Decimal(summary["show_price"])
		summary["inner_price"] = util.ConvertAmount2Decimal(summary["inner_price"])
		summary["recycling_price"] = util.ConvertAmount2Decimal(summary["recycling_price"])
		summary["amount"] = util.ConvertAmount2Decimal(summary["amount"])
	}

	for _, item := range data {
		result = append(result, &DeliveryOrder{
			Date:           item.Date,
			UserID:         item.UserID,
			UserName:       item.UserName,
			ShowPrice:      util.ConvertAmount2Decimal(item.ShowPrice),
			InnerPrice:     util.ConvertAmount2Decimal(item.InnerPrice),
			RecyclingPrice: util.ConvertAmount2Decimal(item.RecyclingPrice),
		})
	}

	return
}

func Format2Excel(dateRange [2]time.Time, _data []*dao.DeliveryOrder) (excelModel *excel.Excel[*DeliveryOrder], err error) {
	_, data := Format(dateRange, nil, _data)

	reflectMap := map[string]func(source *DeliveryOrder) any{
		"日期":   func(source *DeliveryOrder) any { return source.Date },
		"用户ID": func(source *DeliveryOrder) any { return source.UserID },
		"用户昵称": func(source *DeliveryOrder) any { return source.UserName },
		"展示价":  func(source *DeliveryOrder) any { return source.ShowPrice },
		"成本价":  func(source *DeliveryOrder) any { return source.InnerPrice },
		"回收价":  func(source *DeliveryOrder) any { return source.RecyclingPrice },
	}

	excelModel = &excel.Excel[*DeliveryOrder]{
		FileName:   fmt.Sprintf("DeliveryOrder_%s-%s", dateRange[0].Format(pkg.FILE_DATE_FORMAT), dateRange[1].Format(pkg.FILE_DATE_FORMAT)),
		SheetNames: []string{"发货报表"},
		SheetNameWithHead: map[string][]string{
			"发货报表": {
				"日期",
				"用户ID", "用户昵称", "展示价", "成本价",
				"回收价",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*DeliveryOrder]{
			"发货报表": data,
		},
		ReflectMap: map[string]excel.RowReflect[*DeliveryOrder]{
			"发货报表": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
