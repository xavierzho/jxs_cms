package form

import (
	"time"

	"data_backend/apps/v2/internal/report/revenue/dao"
	"data_backend/pkg"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

// 参与数据
type Pating struct {
	Date       string          `json:"date"`
	UserCnt    uint            `json:"user_cnt"` // 参与用户数: 参与各类潮玩+集市(创建, 下单)
	Rate       decimal.Decimal `json:"rate"`     // 参与用户/日活
	UserCntNew uint            `json:"user_cnt_new"`
	RateNew    decimal.Decimal `json:"rate_new"` // 参与用户(新)/日活(新)
}

func FormatPating(dateRange [2]time.Time, patingData []*dao.Pating, activeData []*dao.Active) (result []Pating, err error) {
	var patingDataMap = make(map[string]dao.Pating, len(patingData))
	for _, item := range patingData {
		patingDataMap[item.Date] = *item
	}
	var activeDataMap = make(map[string]dao.Active, len(activeData))
	for _, item := range activeData {
		activeDataMap[item.Date] = *item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)
		item := Pating{
			Date:       cDateStr,
			UserCnt:    patingDataMap[cDateStr].UserCnt,
			UserCntNew: patingDataMap[cDateStr].UserCntNew,
		}

		item.Rate = util.SaveRatio2Decimal(item.UserCnt, activeDataMap[item.Date].ActiveCnt)
		item.RateNew = util.SaveRatio2Decimal(item.UserCntNew, activeDataMap[item.Date].ActiveCntNew)

		result = append(result, item)
	}

	return result, nil
}
