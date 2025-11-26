package form

import (
	"context"
	"fmt"
	"strconv"

	"data_backend/apps/v2/internal/inquire/gacha/dao"
	"data_backend/internal/app"
	"data_backend/internal/global"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type RevenueRequest struct {
	*app.Pager
	IsBoxDim         bool      `form:"is_box_dim"`           // branch 是否按 box 统计
	TypeList         []int     `form:"type_list[]"`          // GMParams
	GachaName        string    `form:"gacha_name"`           // GMParams
	BetRate          *[2]int   `form:"bet_rate[]"`           // AwardParams
	RevenueRange     *[2]int64 `form:"revenue_range[]"`      // OutParams
	RevenueRateRange *[2]int   `form:"revenue_rate_range[]"` // OutParams
}

func (q *RevenueRequest) Parse() (paramsGroup dao.RevenueRequestParamsGroup, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	q.Pager.Parse()

	paramsGroup.IsBoxDim = q.IsBoxDim

	// GMParams
	{
		if len(q.TypeList) != 0 {
			typeList := []int{}
			for _, i := range q.TypeList {
				switch i {
				case 101, 102, 103, 104, 105, 106:
					typeList = append(typeList, i%100)
				}
			}

			if len(typeList) != 0 {
				paramsGroup.GMParams = append(paramsGroup.GMParams, database.QueryWhere{
					Prefix: "gm.type in ?",
					Value:  []any{typeList},
				})
			}
		}
		if q.GachaName != "" {
			paramsGroup.GMParams = append(paramsGroup.GMParams, database.QueryWhere{
				Prefix: "gm.name = ?",
				Value:  []any{q.GachaName},
			})
		}
	}

	// AwardParams
	{
		if q.BetRate != nil {
			paramsGroup.AwardParams = append(paramsGroup.AwardParams, database.QueryWhere{
				Prefix: "(bet_nums/total_nums) between ? and ?",
				Value:  []any{util.Divide2Float64(q.BetRate[0], 100), util.Divide2Float64(q.BetRate[1], 100)},
			})
		}
	}

	// OutParams
	{
		if q.RevenueRange != nil {
			paramsGroup.OutParams = append(paramsGroup.OutParams, database.QueryWhere{
				Prefix: "(ifnull(amount, 0) - inner_price_bet) between ? and ?",
				Value:  []any{util.ReconvertAmount2Decimal(q.RevenueRange[0]).IntPart(), util.ReconvertAmount2Decimal(q.RevenueRange[1]).IntPart()},
			})
		}

		if q.RevenueRateRange != nil {
			paramsGroup.OutParams = append(paramsGroup.OutParams, database.QueryWhere{
				Prefix: "if(amount is null, 0, (amount - inner_price_bet)/amount) between ? and ?",
				Value:  []any{util.Divide2Float64(q.RevenueRateRange[0], 100), util.Divide2Float64(q.RevenueRateRange[1], 100)},
			})
		}
	}

	return
}

func (q *RevenueRequest) Valid() (err error) {
	for _, i := range q.TypeList {
		switch i {
		case 101, 102, 103, 104, 105, 106:
		default:
			return fmt.Errorf("not expected type: %d", q.TypeList)
		}
	}

	if q.BetRate != nil && (q.BetRate[1] < q.BetRate[0] || q.BetRate[0] < 0 || q.BetRate[1] > 100) {
		return fmt.Errorf("invalid BetRate: %v", q.BetRate)
	}

	if q.RevenueRange != nil && q.RevenueRange[1] < q.RevenueRange[0] {
		return fmt.Errorf("invalid RevenueRange: %v", q.RevenueRange)
	}

	if q.RevenueRateRange != nil && q.RevenueRateRange[1] < q.RevenueRateRange[0] {
		return fmt.Errorf("invalid RevenueRateRange: %v", q.RevenueRateRange)
	}

	return nil
}

type GachaRevenue struct {
	GachaID              string          `json:"gacha_id"`
	GachaType            int             `json:"gacha_type"`
	GachaTypeStr         string          `json:"gacha_type_str"`
	GachaName            string          `json:"gacha_name"`
	BoxOutNo             int64           `json:"box_out_no"`
	BetNums              int             `json:"bet_nums"`
	TotalNums            int             `json:"total_nums"`
	BetRate              float64         `json:"bet_rate"`
	Price                decimal.Decimal `json:"price"`
	DiscountPrice        decimal.Decimal `json:"discount_price"`
	Amount               decimal.Decimal `json:"amount"`
	AmountLeft           decimal.Decimal `json:"amount_left"`
	InnerPriceBet        decimal.Decimal `json:"inner_price_bet"`
	InnerPriceLeft       decimal.Decimal `json:"inner_price_left"`
	InnerPriceLeftNormal decimal.Decimal `json:"inner_price_left_normal"`
	InnerPriceBetExtra   decimal.Decimal `json:"inner_price_bet_extra"`
	InnerPriceLeftExtra  decimal.Decimal `json:"inner_price_left_extra"`
	Revenue              decimal.Decimal `json:"revenue"`
	RevenueRate          decimal.Decimal `json:"revenue_rate"`
	SPLeftNum            int64           `json:"sp_left_num"`
}

func FormatRevenue(ctx context.Context, _summary map[string]any, data []*dao.GachaRevenue) (summary map[string]any, result []*GachaRevenue) {
	summary = _summary
	summary["amount"] = util.ConvertAmount2Decimal(summary["amount"])
	summary["amount_left"] = util.ConvertAmount2Decimal(summary["amount_left"])
	summary["inner_price_bet"] = util.ConvertAmount2Decimal(summary["inner_price_bet"])
	summary["inner_price_left"] = util.ConvertAmount2Decimal(summary["inner_price_left"])
	summary["inner_price_left_normal"] = util.ConvertAmount2Decimal(summary["inner_price_left_normal"])
	summary["inner_price_bet_extra"] = util.ConvertAmount2Decimal(summary["inner_price_bet_extra"])
	summary["inner_price_left_extra"] = util.ConvertAmount2Decimal(summary["inner_price_left_extra"])
	summary["revenue"] = util.Sub2Decimal(summary["amount"], summary["inner_price_bet"])
	summary["revenue_rate"] = util.SaveRatio2Decimal(summary["revenue"], summary["amount"])

	for _, item := range data {
		var gachaName = item.GachaName
		if item.Period != 0 {
			gachaName = fmt.Sprintf("%s 第%d期", item.GachaName, item.Period)
		}

		resultItem := &GachaRevenue{
			GachaID:              strconv.FormatInt(item.GachaID, 10),
			GachaType:            item.GachaType,
			GachaTypeStr:         global.I18n.T(ctx, "source_type", convert.GetString(item.GachaType)),
			GachaName:            gachaName,
			BoxOutNo:             item.BoxOutNo,
			BetNums:              item.BetNums,
			TotalNums:            item.TotalNums,
			BetRate:              util.SaveRatio2Float64(item.BetNums, item.TotalNums),
			Price:                util.ConvertAmount2Decimal(item.Price),
			DiscountPrice:        util.ConvertAmount2Decimal(item.DiscountPrice),
			Amount:               util.ConvertAmount2Decimal(item.Amount),
			InnerPriceBet:        util.ConvertAmount2Decimal(item.InnerPriceBetNormal + item.InnerPriceBetExtra),
			InnerPriceLeft:       util.ConvertAmount2Decimal(item.InnerPriceLeftNormal + item.InnerPriceLeftExtra),
			InnerPriceLeftNormal: util.ConvertAmount2Decimal(item.InnerPriceLeftNormal),
			InnerPriceBetExtra:   util.ConvertAmount2Decimal(item.InnerPriceBetExtra),
			InnerPriceLeftExtra:  util.ConvertAmount2Decimal(item.InnerPriceLeftExtra),
			SPLeftNum:            item.SPLeftNum,
		}

		leftNums := resultItem.TotalNums - resultItem.BetNums
		if resultItem.DiscountPrice.IsZero() {
			resultItem.AmountLeft = resultItem.Price.Mul(decimal.NewFromInt(int64(leftNums)))
		} else {
			resultItem.AmountLeft = resultItem.DiscountPrice.Mul(decimal.NewFromInt(int64(leftNums)))
		}

		resultItem.Revenue = util.Sub2Decimal(resultItem.Amount, resultItem.InnerPriceBet)
		resultItem.RevenueRate = util.SaveRatio2Decimal(resultItem.Revenue, resultItem.Amount)

		result = append(result, resultItem)
	}

	return
}
