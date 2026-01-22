package form

import (
	"fmt"
	"time"

	"data_backend/apps/v2/internal/report/market/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
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

func (q *ListRequest) Parse() (dateRange [2]time.Time, err error) {
	q.Pager.Parse()

	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateRangeRequest
}

type Market struct {
	Date     string          `json:"date"`
	UserCnt  uint            `json:"user_cnt"`
	OrderCnt uint            `json:"order_cnt"`
	Amount   decimal.Decimal `json:"amount"`
	Amount0  decimal.Decimal `json:"amount_0"`
	Amount1  decimal.Decimal `json:"amount_1"`
	Amount2  decimal.Decimal `json:"amount_2"`
}

func Format(dateRange [2]time.Time, _summary map[string]any, data []*dao.Market) (summary map[string]any, result []Market, err error) {
	summary = _summary
	if summary != nil {
		summary["amount_0"] = util.ConvertAmount2Decimal(summary["amount_0"])
		summary["amount_1"] = util.ConvertAmount2Decimal(summary["amount_1"])
		summary["amount_2"] = util.ConvertAmount2Decimal(summary["amount_2"])
		summary["amount"] = summary["amount_0"].(decimal.Decimal).Add(summary["amount_1"].(decimal.Decimal)).Add(summary["amount_2"].(decimal.Decimal))
	}

	var dataMap = make(map[string]dao.Market)
	for _, item := range data {
		dataMap[item.Date] = *item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)
		item := Market{
			Date:     cDateStr,
			UserCnt:  dataMap[cDateStr].UserCnt,
			OrderCnt: dataMap[cDateStr].OrderCnt,
			Amount0:  util.ConvertAmount2Decimal(dataMap[cDateStr].Amount0),
			Amount1:  util.ConvertAmount2Decimal(dataMap[cDateStr].Amount1),
			Amount2:  util.ConvertAmount2Decimal(dataMap[cDateStr].Amount2),
		}

		item.Amount = item.Amount0.Add(item.Amount1).Add(item.Amount2)

		result = append(result, item)
	}

	return
}

func Format2Excel(dateRange [2]time.Time, _data []*dao.Market) (excelModel *excel.Excel[Market], err error) {
	_, data, err := Format(dateRange, nil, _data)
	if err != nil {
		return nil, err
	}

	reflectMap := map[string]func(source Market) any{
		"日期":      func(source Market) any { return source.Date },
		"参与用户数":   func(source Market) any { return source.UserCnt },
		"新增订单数":   func(source Market) any { return source.OrderCnt },
		"成交金额(总)": func(source Market) any { return source.Amount },
		"成交金额0":   func(source Market) any { return source.Amount0 },
		"成交金额1":   func(source Market) any { return source.Amount1 },
		"成交金额2":   func(source Market) any { return source.Amount2 },
	}

	excelModel = &excel.Excel[Market]{
		FileName:   fmt.Sprintf("market_%s-%s", dateRange[0].Format(pkg.FILE_DATE_FORMAT), dateRange[1].Format(pkg.FILE_DATE_FORMAT)),
		SheetNames: []string{"集市报表"},
		SheetNameWithHead: map[string][]string{
			"集市报表": {
				"日期", "参与用户数", "新增订单数",
				"成交金额(总)", "成交金额0", "成交金额1", "成交金额2",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[Market]{
			"集市报表": data,
		},
		ReflectMap: map[string]excel.RowReflect[Market]{
			"集市报表": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
