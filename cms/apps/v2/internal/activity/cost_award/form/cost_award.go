package form

import (
	"fmt"
	"time"

	"data_backend/apps/v2/internal/activity/cost_award/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

var COST_AWARD_POINT_STEP = decimal.NewFromInt(10)

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

type CostAward struct {
	Date                string          `json:"date"`
	GetUserCnt          uint            `json:"get_user_cnt"`
	GetPoint            decimal.Decimal `json:"get_point"`
	AcceptUserCnt       uint            `json:"accept_user_cnt"`
	AcceptPoint         decimal.Decimal `json:"accept_point"`
	AwardAmount         decimal.Decimal `json:"award_amount"`
	AwardItemShowPrice  decimal.Decimal `json:"award_item_show_price"`
	AwardItemInnerPrice decimal.Decimal `json:"award_item_inner_price"`
}

func Format(dateRange [2]time.Time, _summary map[string]any, data []*dao.CostAward) (summary map[string]any, result []*CostAward) {
	summary = _summary
	if summary != nil {
		summary["get_point"] = util.ConvertAmount2Decimal(summary["get_amount"]).Mul(COST_AWARD_POINT_STEP)
		summary["accept_point"] = util.ConvertAmount2Decimal(summary["accept_amount"]).Mul(COST_AWARD_POINT_STEP)
		summary["award_amount"] = util.ConvertAmount2Decimal(summary["award_amount"])
		summary["award_item_show_price"] = util.ConvertAmount2Decimal(summary["award_item_show_price"])
		summary["award_item_inner_price"] = util.ConvertAmount2Decimal(summary["award_item_inner_price"])
	}

	var dataMap = make(map[string]dao.CostAward)
	for _, item := range data {
		dataMap[item.Date] = *item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)
		result = append(result, &CostAward{
			Date:                cDateStr,
			GetUserCnt:          dataMap[cDateStr].GetUserCnt,
			GetPoint:            util.ConvertAmount2Decimal(dataMap[cDateStr].GetAmount).Mul(COST_AWARD_POINT_STEP),
			AcceptUserCnt:       dataMap[cDateStr].AcceptUserCnt,
			AcceptPoint:         util.ConvertAmount2Decimal(dataMap[cDateStr].AcceptAmount).Mul(COST_AWARD_POINT_STEP),
			AwardAmount:         util.ConvertAmount2Decimal(dataMap[cDateStr].AwardAmount),
			AwardItemShowPrice:  util.ConvertAmount2Decimal(dataMap[cDateStr].AwardItemShowPrice),
			AwardItemInnerPrice: util.ConvertAmount2Decimal(dataMap[cDateStr].AwardItemInnerPrice),
		})
	}

	return
}

func Format2Excel(dateRange [2]time.Time, _data []*dao.CostAward) (excelModel *excel.Excel[*CostAward], err error) {
	_, data := Format(dateRange, nil, _data)

	reflectMap := map[string]func(source *CostAward) any{
		"日期":        func(source *CostAward) any { return source.Date },
		"获得用户数":     func(source *CostAward) any { return source.GetUserCnt },
		"获得总额":      func(source *CostAward) any { return source.GetPoint },
		"领取用户数":     func(source *CostAward) any { return source.AcceptUserCnt },
		"领取总额":      func(source *CostAward) any { return source.AcceptPoint },
		"现金奖励总额":    func(source *CostAward) any { return source.AwardAmount },
		"物品奖励展示价总额": func(source *CostAward) any { return source.AwardItemShowPrice },
		"物品奖励成本价总额": func(source *CostAward) any { return source.AwardItemInnerPrice },
	}

	excelModel = &excel.Excel[*CostAward]{
		FileName:   fmt.Sprintf("cost_award_%s-%s", dateRange[0].Format(pkg.FILE_DATE_FORMAT), dateRange[1].Format(pkg.FILE_DATE_FORMAT)),
		SheetNames: []string{"欧气值报表"},
		SheetNameWithHead: map[string][]string{
			"欧气值报表": {
				"日期",
				"获得用户数", "获得总额", "领取用户数", "领取总额",
				"现金奖励总额", "物品奖励展示价总额", "物品奖励成本价总额",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*CostAward]{
			"欧气值报表": data,
		},
		ReflectMap: map[string]excel.RowReflect[*CostAward]{
			"欧气值报表": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
