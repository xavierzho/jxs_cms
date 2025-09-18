package form

import (
	"time"

	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

// 营收数据
// 单位: 元
type Summary struct {
	Date                       string          `json:"date"`
	ActiveCnt                  uint            `json:"active_cnt"`                    // 日活;活跃用户数;登录用户数
	WalletBalance              decimal.Decimal `json:"wallet_balance"`                // 钱包余额
	MerchantBalance            decimal.Decimal `json:"merchant_balance"`              // 商户钱包余额
	PayAmount                  decimal.Decimal `json:"pay_amount"`                    // 总付费
	PayAmountBet               decimal.Decimal `json:"pay_amount_bet"`                // 总付费(抽赏)
	RechargeAmount             decimal.Decimal `json:"recharge_amount"`               // 充值金额
	RechargeAmountWeChat       decimal.Decimal `json:"recharge_amount_wechat"`        // 充值金额 微信
	RechargeAmountAli          decimal.Decimal `json:"recharge_amount_ali"`           // 充值金额 支付宝
	RechargeRefundAmount       decimal.Decimal `json:"recharge_refund_amount"`        // 充值退款金额
	RechargeRefundAmountWeChat decimal.Decimal `json:"recharge_refund_amount_wechat"` // 充值退款金额 微信
	RechargeRefundAmountAli    decimal.Decimal `json:"recharge_refund_amount_ali"`    // 充值退款金额 支付宝
	DrawAmount                 decimal.Decimal `json:"draw_amount"`                   // 总退款(￥)
	TaxAmount                  decimal.Decimal `json:"tax_amount"`                    // 总抽水
	Revenue                    decimal.Decimal `json:"revenue"`                       // 总营收: 日付费-日退款(￥)
	RevenueRate                decimal.Decimal `json:"revenue_rate"`                  // 营收率: 浮动营收/日付费
	RevenueARPU                decimal.Decimal `json:"revenue_arpu"`                  // 营收ARPU: 浮动营收/日活
	RefundAmount               decimal.Decimal `json:"refund_amount"`                 // 总退款(潮币)

}

func FormatSummary(dateRange [2]time.Time, data []map[string]interface{}) (result []Summary, err error) {
	var dataMap = make(map[string]map[string]interface{}, len(data))
	for _, item := range data {
		dataMap[item["date"].(string)] = item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)

		dataItem := dataMap[cDateStr]
		if dataItem == nil {
			dataItem = make(map[string]interface{})
		}

		item := Summary{
			Date:                       cDateStr,
			ActiveCnt:                  convert.GetUint(dataItem["active_cnt"]),
			WalletBalance:              util.ConvertAmount2Decimal(dataItem["wallet_balance"]),
			MerchantBalance:            util.ConvertAmount2Decimal(dataItem["merchant_balance"]),
			PayAmount:                  util.ConvertAmount2Decimal(dataItem["pay_amount"]),
			PayAmountBet:               util.ConvertAmount2Decimal(dataItem["pay_amount_bet"]),
			RechargeAmount:             util.ConvertAmount2Decimal(dataItem["recharge_amount"]),
			RechargeAmountWeChat:       util.ConvertAmount2Decimal(dataItem["recharge_amount_wechat"]),
			RechargeAmountAli:          util.ConvertAmount2Decimal(dataItem["recharge_amount_ali"]),
			RechargeRefundAmount:       util.ConvertAmount2Decimal(dataItem["recharge_refund_amount"]),
			RechargeRefundAmountWeChat: util.ConvertAmount2Decimal(dataItem["recharge_refund_amount_wechat"]),
			RechargeRefundAmountAli:    util.ConvertAmount2Decimal(dataItem["recharge_refund_amount_ali"]),
			DrawAmount:                 util.ConvertAmount2Decimal(dataItem["draw_amount"]),
			TaxAmount:                  util.ConvertAmount2Decimal(dataItem["tax_amount"]),
			RefundAmount:               util.ConvertAmount2Decimal(dataItem["refund_amount"]),
		}

		item.Revenue = item.RechargeAmount.Sub(item.RechargeRefundAmount).Sub(item.DrawAmount)
		item.RevenueRate = util.SaveRatio2Decimal(item.Revenue, item.RechargeAmount)
		item.RevenueARPU = util.SaveDivide2Decimal(item.Revenue, item.ActiveCnt)

		result = append(result, item)
	}

	return result, nil
}
