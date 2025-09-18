package form

import (
	"time"

	"data_backend/apps/v2/internal/report/revenue/dao"
	"data_backend/pkg"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type Active struct {
	Date            string          `json:"date"`
	ActivateCnt     uint            `json:"activate_cnt"`      // 打开设备数
	ActivateCntNew  uint            `json:"activate_cnt_new"`  // 非注册用户打开app设备数
	RegisterCnt     uint            `json:"register_cnt"`      // 新注册用户数
	RegisterCntRate decimal.Decimal `json:"register_cnt_rate"` // 注册率: RegisterCnt/ActivateCntNew
	ActiveCnt       uint            `json:"active_cnt"`        // 日活;活跃用户数;登录用户数
	ActiveCntNew    uint            `json:"active_cnt_new"`
	ActiveCntOld    uint            `json:"active_cnt_old"`
	MaxOnlineCnt    uint            `json:"max_online_cnt"`  // 最大在线用户数
	ValidatedCnt7   uint            `json:"validated_cnt_7"` // 有效用户数，7天内登录两次以上且消费过
	ValidatedCnt15  uint            `json:"validated_cnt_15"`
	PerOnlineTime   decimal.Decimal `json:"per_online_time"`
	WastageRate1    decimal.Decimal `json:"wastage_rate_1"` // 流失率 2-7天均未登录
	WastageRate3    decimal.Decimal `json:"wastage_rate_3"` // 流失率 1-3天登录过, 4-7天均未登录
}

func FormatActive(dateRange [2]time.Time, activeData []*dao.Active, wastageData []*dao.Wastage) (result []Active, err error) {
	var activeDataMap = make(map[string]dao.Active, len(activeData))
	for _, item := range activeData {
		activeDataMap[item.Date] = *item
	}
	var wastageDataMap = make(map[string]dao.Wastage, len(wastageData))
	for _, item := range wastageData {
		wastageDataMap[item.Date] = *item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)
		item := Active{
			Date:           cDateStr,
			ActivateCnt:    activeDataMap[cDateStr].ActivateCnt,
			ActivateCntNew: activeDataMap[cDateStr].ActivateCntNew,
			RegisterCnt:    activeDataMap[cDateStr].RegisterCnt,
			ActiveCnt:      activeDataMap[cDateStr].ActiveCnt,
			ActiveCntNew:   activeDataMap[cDateStr].ActiveCntNew,
			ActiveCntOld:   activeDataMap[cDateStr].ActiveCntOld,
			MaxOnlineCnt:   activeDataMap[cDateStr].MaxOnlineCnt,
			ValidatedCnt7:  activeDataMap[cDateStr].ValidatedCnt7,
			ValidatedCnt15: activeDataMap[cDateStr].ValidatedCnt15,
		}

		item.RegisterCntRate = util.SaveRatio2Decimal(item.RegisterCnt, item.ActivateCntNew)
		item.WastageRate1 = util.SaveRatio2Decimal(wastageDataMap[item.Date].WastageCnt1, item.RegisterCnt)
		item.WastageRate3 = util.SaveRatio2Decimal(wastageDataMap[item.Date].WastageCnt3, item.RegisterCnt)

		result = append(result, item)
	}

	return result, nil
}
