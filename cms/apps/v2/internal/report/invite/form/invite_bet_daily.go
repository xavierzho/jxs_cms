package form

import (
	"fmt"
	"time"

	"data_backend/apps/v2/internal/report/invite/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type GenerateDailyRequest struct {
	iForm.DateRangeRequest
}

type ListDailyRequest struct {
	app.Pager
	AllDailyRequest
}

func (q *ListDailyRequest) Parse() (dateRange [2]time.Time, err error) {
	q.Pager.Parse()

	return q.AllDailyRequest.Parse()
}

type AllDailyRequest struct {
	iForm.DateRangeRequest
}

type InviteBetDaily struct {
	Date        string          `json:"date"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Amount      decimal.Decimal `json:"amount"`
	Difference  decimal.Decimal `json:"difference"`
}

func FormatDaily(dateRange [2]time.Time, _summary map[string]any, data []*dao.InviteBetDaily) (summary map[string]any, result []*InviteBetDaily) {
	summary = _summary
	if summary != nil {
		summary["amount"] = util.ConvertAmount2Decimal(summary["amount"])
		summary["total_amount"] = util.ConvertAmount2Decimal(summary["total_amount"])
		summary["difference"] = util.ConvertAmount2Decimal(summary["difference"])
	}

	var dataMap = make(map[string]dao.InviteBetDaily)
	for _, item := range data {
		dataMap[item.Date] = *item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)
		result = append(result, &InviteBetDaily{
			Date:        cDateStr,
			TotalAmount: util.ConvertAmount2Decimal(dataMap[cDateStr].TotalAmount),
			Amount:      util.ConvertAmount2Decimal(dataMap[cDateStr].Amount),
			Difference:  util.ConvertAmount2Decimal(dataMap[cDateStr].Difference),
		})
	}

	return
}

func Format2DailyExcel(dateRange [2]time.Time, _data []*dao.InviteBetDaily) (excelModel *excel.Excel[*InviteBetDaily], err error) {
	_, data := FormatDaily(dateRange, nil, _data)

	reflectMap := map[string]func(source *InviteBetDaily) any{
		"日期":       func(source *InviteBetDaily) any { return source.Date },
		"总抽赏":      func(source *InviteBetDaily) any { return source.TotalAmount },
		"被邀用户抽赏金额": func(source *InviteBetDaily) any { return source.Amount },
		"差额":       func(source *InviteBetDaily) any { return source.Difference },
	}

	excelModel = &excel.Excel[*InviteBetDaily]{
		FileName:   fmt.Sprintf("invite_bet_daily_%s-%s", dateRange[0].Format(pkg.FILE_DATE_FORMAT), dateRange[1].Format(pkg.FILE_DATE_FORMAT)),
		SheetNames: []string{"邀请日报表"},
		SheetNameWithHead: map[string][]string{
			"邀请日报表": {
				"日期",
				"总抽赏", "被邀用户抽赏金额", "差额",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*InviteBetDaily]{
			"邀请日报表": data,
		},
		ReflectMap: map[string]excel.RowReflect[*InviteBetDaily]{
			"邀请日报表": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
