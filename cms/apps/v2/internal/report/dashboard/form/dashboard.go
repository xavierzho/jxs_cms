package form

import (
	"time"

	"data_backend/apps/v2/internal/report/dashboard/dao"
	"data_backend/pkg"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type GenerateRequest struct {
	DateRange [2]string `form:"date_range[]" binding:"required"`
}

func (q *GenerateRequest) Parse() (dateRange [2]time.Time, err error) {
	if dateRange, err = util.ParseDateParams(q.DateRange); err != nil {
		return
	}

	return
}

type Dashboard struct {
	Date                 string          `json:"date"`
	NewUserCnt           int             `json:"new_user_cnt"`
	ActiveUserCnt        int             `json:"active_user_cnt"`
	PatingUserCnt        int             `json:"pating_user_cnt"`
	PatingUserCntNew     int             `json:"pating_user_cnt_new"`
	PatingRateNew        decimal.Decimal `json:"pating_rate_new"`
	PayUserCnt           int             `json:"pay_user_cnt"`
	PayUserCntNew        int             `json:"pay_user_cnt_new"`
	RechargeUserCnt      int             `json:"recharge_user_cnt"`
	RechargeUserCntNew   int             `json:"recharge_user_cnt_new"`
	RechargeAmount       decimal.Decimal `json:"recharge_amount"`
	RechargeAmountWeChat decimal.Decimal `json:"recharge_amount_wechat"`
	RechargeAmountAli    decimal.Decimal `json:"recharge_amount_ali"`
	DrawAmount           decimal.Decimal `json:"draw_amount"`
}

// 填充
func Format(dateRange [2]time.Time, data ...*dao.Dashboard) (result []*Dashboard) {
	var dataMap = make(map[string]*dao.Dashboard, len(data))
	for _, item := range data {
		dataMap[item.Date] = item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)

		dataItem := dataMap[cDateStr]
		if dataItem == nil {
			dataItem = &dao.Dashboard{}
		}

		result = append(result, &Dashboard{
			Date:                 cDateStr,
			NewUserCnt:           dataItem.NewUserCnt,
			ActiveUserCnt:        dataItem.ActiveUserCnt,
			PatingUserCnt:        dataItem.PatingUserCnt,
			PatingUserCntNew:     dataItem.PatingUserCntNew,
			PatingRateNew:        util.SaveRatio2Decimal(dataItem.PatingUserCntNew, dataItem.NewUserCnt),
			PayUserCnt:           dataItem.PayUserCnt,
			PayUserCntNew:        dataItem.PayUserCntNew,
			RechargeUserCnt:      dataItem.RechargeUserCnt,
			RechargeUserCntNew:   dataItem.RechargeUserCntNew,
			RechargeAmount:       util.ConvertAmount2Decimal(dataItem.RechargeAmount),
			RechargeAmountWeChat: util.ConvertAmount2Decimal(dataItem.RechargeAmountWeChat),
			RechargeAmountAli:    util.ConvertAmount2Decimal(dataItem.RechargeAmountAli),
			DrawAmount:           util.ConvertAmount2Decimal(dataItem.DrawAmount),
		})
	}

	return
}
