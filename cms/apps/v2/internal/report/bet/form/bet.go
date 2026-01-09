package form

import (
	"fmt"
	"time"

	"data_backend/apps/v2/internal/report/bet/dao"
	"data_backend/pkg"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

const (
	BET_TYPE_FIRST_PRIZE = "FirstPrize" // 一番赏 // gacha_type = 1 source_type = 101
	BET_TYPE_GASHAPON    = "Gashapon"   // 扭蛋机 // gacha_type = 2 source_type = 102
	BET_TYPE_CHAO        = "Chao"       // 潮玩赏 // gacha_type = 3 source_type = 103
	BET_TYPE_HOLE        = "Hole"       // 洞洞乐 // gacha_type = 4 source_type = 104
	BET_TYPE_CLUSTER     = "ClusterBox" // 擂台赏 // gacha_type = 5 source_type = 105
	BET_TYPE_CHUANCHUAN  = "ChuanChuan" // 串串赏 // gacha_type = 6 source_type = 106
)

var BET_TYPE_LIST = []string{BET_TYPE_FIRST_PRIZE, BET_TYPE_GASHAPON, BET_TYPE_CHAO, BET_TYPE_HOLE, BET_TYPE_CLUSTER, BET_TYPE_CHUANCHUAN}

type GenerateRequest struct {
	DateRange [2]string `form:"date_range[]" binding:"required"`
}

func (q *GenerateRequest) Parse() (dateRange [2]time.Time, err error) {
	if dateRange, err = util.ParseDateParams(q.DateRange); err != nil {
		return dateRange, err
	}

	return dateRange, nil
}

type AllRequest struct {
	DateRange [2]string `form:"date_range[]" binding:"required"`
	DataType  string    `form:"data_type" binding:"required"`
}

func (q *AllRequest) Parse() (dateRange [2]time.Time, err error) {
	if err = q.Valid(); err != nil {
		return dateRange, err
	}

	if dateRange, err = util.ParseDateParams(q.DateRange); err != nil {
		return dateRange, err
	}

	return dateRange, nil
}

func (q *AllRequest) Valid() (err error) {
	switch q.DataType {
	case BET_TYPE_FIRST_PRIZE, BET_TYPE_GASHAPON, BET_TYPE_CHAO, BET_TYPE_HOLE, BET_TYPE_CLUSTER, BET_TYPE_CHUANCHUAN:
	default:
		return fmt.Errorf("not expected data_type: " + q.DataType)
	}

	return nil
}

type Bet struct {
	Date            string          `json:"date"`
	UserCnt         uint            `json:"user_cnt"`
	BetNums         uint            `json:"bet_nums"`
	BoxCntRemaining uint            `json:"box_cnt_remaining"`
	BoxCntNew       uint            `json:"box_cnt_new"`
	BoxCntClose     uint            `json:"box_cnt_close"`
	Amount          decimal.Decimal `json:"amount"`
	AmountBalance   decimal.Decimal `json:"amount_balance"`
	AmountWeChat    decimal.Decimal `json:"amount_wechat"`
	AmountAli       decimal.Decimal `json:"amount_ali"`
	AmountHuiFu     decimal.Decimal `json:"amount_huifu"`
}

func Format(dateRange [2]time.Time, data []*dao.Bet) (result []Bet, err error) {
	var dataMap = make(map[string]dao.Bet)
	for _, item := range data {
		dataMap[item.Date] = *item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)
		item := Bet{
			Date:            cDateStr,
			UserCnt:         dataMap[cDateStr].UserCnt,
			BetNums:         dataMap[cDateStr].BetNums,
			BoxCntRemaining: dataMap[cDateStr].BoxCntRemaining,
			BoxCntNew:       dataMap[cDateStr].BoxCntNew,
			BoxCntClose:     dataMap[cDateStr].BoxCntClose,
			Amount:          util.ConvertAmount2Decimal(dataMap[cDateStr].Amount),
			AmountWeChat:    util.ConvertAmount2Decimal(dataMap[cDateStr].AmountWeChat),
			AmountAli:       util.ConvertAmount2Decimal(dataMap[cDateStr].AmountAli),
			AmountHuiFu:     util.ConvertAmount2Decimal(dataMap[cDateStr].AmountHuiFu),
		}

		result = append(result, item)
	}

	return result, nil
}
