package form

import (
	"time"

	"data_backend/apps/v2/internal/report/revenue/dao"
	"data_backend/pkg"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

// 付费数据
// 单位: 元
type Pay struct {
	Date            string          `json:"date"`
	Amount          decimal.Decimal `json:"amount"` // 总消费金额: 潮玩+集市
	AmountNew       decimal.Decimal `json:"amount_new"`
	AmountOld       decimal.Decimal `json:"amount_old"`
	UserCnt         uint            `json:"user_cnt"` // 总消费人数
	UserCntNew      uint            `json:"user_cnt_new"`
	UserCntOld      uint            `json:"user_cnt_old"`
	UserCntFirst    uint            `json:"user_cnt_first"`   // 首次消费人数
	ARPU            decimal.Decimal `json:"arpu"`             // 日付费(总)/日活
	ARPPU           decimal.Decimal `json:"arppu"`            // 总付费/付费人数(总)
	ARPPUNew        decimal.Decimal `json:"arppu_new"`        // 日付费(新注册)/付费用户(新注册)
	ARPPUOld        decimal.Decimal `json:"arppu_old"`        // 日付费(老用户)/付费人数(老用户)
	Permeability    decimal.Decimal `json:"permeability"`     // 渗透率 付费人数(总)/日活
	PermeabilityNew decimal.Decimal `json:"permeability_new"` // 渗透率 付费人数(新注册)/新增注册
	PermeabilityOld decimal.Decimal `json:"permeability_old"` // 渗透率 付费人数(老用户)/老用户日活
}

func FormatPay(dateRange [2]time.Time, payData []*dao.Pay, activeData []*dao.Active) (result []Pay, err error) {
	var payDataMap = make(map[string]dao.Pay, len(payData))
	for _, item := range payData {
		payDataMap[item.Date] = *item
	}
	var activeDataMap = make(map[string]dao.Active, len(activeData))
	for _, item := range activeData {
		activeDataMap[item.Date] = *item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)
		item := Pay{
			Date:         cDateStr,
			Amount:       util.ConvertAmount2Decimal(payDataMap[cDateStr].Amount),
			AmountNew:    util.ConvertAmount2Decimal(payDataMap[cDateStr].AmountNew),
			AmountOld:    util.ConvertAmount2Decimal(payDataMap[cDateStr].AmountOld),
			UserCnt:      payDataMap[cDateStr].UserCnt,
			UserCntNew:   payDataMap[cDateStr].UserCntNew,
			UserCntOld:   payDataMap[cDateStr].UserCntOld,
			UserCntFirst: payDataMap[cDateStr].UserCntFirst,
		}

		item.ARPU = util.SaveDivide2Decimal(item.Amount, activeDataMap[item.Date].ActiveCnt)
		item.ARPPU = util.SaveDivide2Decimal(item.Amount, item.UserCnt)
		item.ARPPUNew = util.SaveDivide2Decimal(item.AmountNew, item.UserCntNew)
		item.ARPPUOld = util.SaveDivide2Decimal(item.AmountOld, item.UserCntOld)

		item.Permeability = util.SaveRatio2Decimal(item.UserCnt, activeDataMap[item.Date].ActiveCnt)
		item.PermeabilityNew = util.SaveRatio2Decimal(item.UserCntNew, activeDataMap[item.Date].ActiveCntNew)
		item.PermeabilityOld = util.SaveRatio2Decimal(item.UserCntOld, activeDataMap[item.Date].ActiveCntOld)

		result = append(result, item)
	}

	return result, nil
}
