package form

import (
	"time"

	"data_backend/apps/v2/internal/report/revenue/dao"
	"data_backend/pkg"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

// 退款(￥)数据
// 单位: 元
type Draw struct {
	Date      string          `json:"date"`
	Amount    decimal.Decimal `json:"amount"`
	UserCnt   uint            `json:"user_cnt"`
	PerAmount decimal.Decimal `json:"per_amount"` // Amount/UserCnt
	Rate      decimal.Decimal `json:"rate"`       // 当日退款(￥)人数/日活
	Tax       decimal.Decimal `json:"tax"`        // 总抽水
	TaxNew    decimal.Decimal `json:"tax_new"`
	TaxOld    decimal.Decimal `json:"tax_old"`
}

func FormatDraw(dateRange [2]time.Time, drawData []*dao.Draw, activeData []*dao.Active) (result []Draw, err error) {
	var drawDataMap = make(map[string]dao.Draw, len(drawData))
	for _, item := range drawData {
		drawDataMap[item.Date] = *item
	}
	var activeDataMap = make(map[string]dao.Active, len(activeData))
	for _, item := range activeData {
		activeDataMap[item.Date] = *item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)
		item := Draw{
			Date:    cDateStr,
			Amount:  util.ConvertAmount2Decimal(drawDataMap[cDateStr].Amount),
			UserCnt: drawDataMap[cDateStr].UserCnt,
		}

		item.PerAmount = util.SaveDivide2Decimal(item.Amount, item.UserCnt)

		item.Rate = util.SaveRatio2Decimal(item.UserCnt, activeDataMap[item.Date].ActiveCnt)

		result = append(result, item)
	}

	return result, nil
}
