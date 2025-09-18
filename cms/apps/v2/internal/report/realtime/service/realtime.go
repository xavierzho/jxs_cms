package service

import (
	"context"
	"sort"
	"sync"
	"time"

	"data_backend/apps/v2/internal/report/realtime/dao"
	"data_backend/apps/v2/internal/report/realtime/form"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"
	"data_backend/pkg/redisdb"
	"data_backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// 实时数据类型
const (
	realtime_data_active            = "active"
	realtime_data_pating            = "pating"
	realtime_data_pay_user_cnt      = "payUserCnt"
	realtime_data_pay_amount        = "payAmount"
	realtime_data_draw_user_cnt     = "drawUserCnt"
	realtime_data_draw_amount       = "drawAmount"
	realtime_data_recharge_user_cnt = "rechargeUserCnt"
	realtime_data_recharge_amount   = "rechargeAmount"
	realtime_data_new_user_cnt      = "newUserCnt"
	realtime_data_online_user_cnt   = "onlineUserCnt"
	realtime_data_summary           = "Summary" // 用于拼接汇总数据
)

// 参与用户类型
const (
	pating_type_first_prize = "FirstPrize" // 一番赏
	pating_type_gashapon    = "Gashapon"   // 扭蛋机
	pating_type_chao        = "Chao"       // 潮玩赏
	pating_type_hole        = "Hole"       // 洞洞乐
	pating_type_market      = "Market"     // 集市
	pating_type_deliver     = "Deliver"    // 发货
	pating_type_shop        = "Shop"       // 商城
)

func realtimeDataRKeyFormatByTime(rKey string, cTime time.Time) string {
	return "realtimeData:" + rKey + ":" + cTime.Format(pkg.DATE_FORMAT)
}

type RealtimeSvc struct {
	ctx    *gin.Context
	rdb    *redisdb.RedisClient
	logger *logger.Logger
	dao    *dao.RealtimeDao
}

func NewRealtimeSvc(ctx *gin.Context, center *gorm.DB, rdb *redisdb.RedisClient, log *logger.Logger) *RealtimeSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".RealtimeSvc")))
	return &RealtimeSvc{
		ctx:    ctx,
		rdb:    rdb,
		logger: log,
		dao:    dao.NewRealtimeDao(center, log),
	}
}

func (svc *RealtimeSvc) Cached(params *form.CachedRequest) (e *errcode.Error) {
	cTime, err := params.Parse()
	if err != nil {
		return errcode.InvalidParams.WithDetails(err.Error())
	}

	endTime := time.Date(cTime.Year(), cTime.Month(), cTime.Day(), cTime.Hour(), (cTime.Minute()/10)*10, 0, 0, pkg.Location)
	startTime := endTime.Add(-10 * time.Minute)
	endTime = endTime.Add(-time.Millisecond)

	eg := errgroup.Group{}
	// 时间段登录用户
	eg.Go(func() error { return svc.cachedActive(startTime, endTime) })

	// 时间段各类型参与用户
	eg.Go(func() error { return svc.cachedPating(startTime, endTime) })

	// 时间段付费金额/人数
	eg.Go(func() error { return svc.cachedPay(startTime, endTime) })

	// 时间段充值金额/人数
	eg.Go(func() error { return svc.cachedRecharge(startTime, endTime) })

	// 时间段退款(￥)金额/人数
	eg.Go(func() error { return svc.cachedDraw(startTime, endTime) })

	// 时间段注册用户数
	eg.Go(func() error { return svc.cachedNewUserCnt(startTime, endTime) })

	// TODO 暂无 onlineUserCnt // 不使用传入的时间而是运行时的时间 // 具体看在线用户如何获取
	eg.Go(func() error { return svc.cachedOnlineUserCnt() })

	err = eg.Wait()
	if err != nil {
		return errcode.ExecuteFail.WithDetails(err.Error())
	}

	return nil
}

// TODO 暂无 登录用户分类/汇总
func (svc *RealtimeSvc) cachedActive(startTime, endTime time.Time) (err error) {
	userCnt, err := svc.dao.GetActiveUserCnt(startTime, endTime)
	if err != nil {
		return err
	}

	return svc.cached(realtime_data_active, endTime, userCnt)
}

func (svc *RealtimeSvc) cachedPating(startTime, endTime time.Time) (err error) {
	// 参与用户数
	data, err := svc.dao.GetPatingUserCnt(startTime, endTime)
	if err != nil {
		return err
	}

	for key, value := range data {
		var patingType = ""
		userCnt := convert.GetInt64(value)
		switch key {
		case "user_cnt":
		case "user_cnt_101":
			patingType = pating_type_first_prize
		case "user_cnt_102":
			patingType = pating_type_gashapon
		case "user_cnt_103":
			patingType = pating_type_chao
		case "user_cnt_104":
			patingType = pating_type_hole
		case "user_cnt_200":
			patingType = pating_type_market
		case "user_cnt_300":
			patingType = pating_type_deliver
		case "user_cnt_600":
			patingType = pating_type_shop
		}

		if err := svc.cached(realtime_data_pating+patingType, endTime, userCnt); err != nil {
			return err
		}
	}

	return nil
}

func (svc *RealtimeSvc) cachedPay(startTime, endTime time.Time) (err error) {
	data, err := svc.dao.GetPayData(startTime, endTime)
	if err != nil {
		return err
	}

	for key, _value := range data {
		var dataType = ""
		var value string
		switch key {
		case "user_cnt":
			dataType = realtime_data_pay_user_cnt
			value = convert.GetString(_value)
		case "amount":
			dataType = realtime_data_pay_amount
			value = util.ConvertAmount2Decimal(_value).String()
		}

		if err := svc.cached(dataType, endTime, value); err != nil {
			return err
		}
	}

	return nil
}

func (svc *RealtimeSvc) cachedRecharge(startTime, endTime time.Time) (err error) {
	data, err := svc.dao.GetRechargeData(startTime, endTime)
	if err != nil {
		return err
	}

	for key, _value := range data {
		var dataType = ""
		var value string
		switch key {
		case "user_cnt":
			dataType = realtime_data_recharge_user_cnt
			value = convert.GetString(_value)
		case "amount":
			dataType = realtime_data_recharge_amount
			value = util.ConvertAmount2Decimal(_value).String()
		}

		if err := svc.cached(dataType, endTime, value); err != nil {
			return err
		}
	}

	return nil
}

func (svc *RealtimeSvc) cachedDraw(startTime, endTime time.Time) (err error) {
	data, err := svc.dao.GetDrawData(startTime, endTime)
	if err != nil {
		return err
	}

	for key, _value := range data {
		var dataType = ""
		var value string
		switch key {
		case "user_cnt":
			dataType = realtime_data_draw_user_cnt
			value = convert.GetString(_value)
		case "amount":
			dataType = realtime_data_draw_amount
			value = util.ConvertAmount2Decimal(_value).String()
		}

		if err := svc.cached(dataType, endTime, value); err != nil {
			return err
		}
	}

	return nil
}

func (svc *RealtimeSvc) cachedNewUserCnt(startTime, endTime time.Time) (err error) {
	userCnt, err := svc.dao.GetNewUserCnt(startTime, endTime)
	if err != nil {
		return err
	}

	return svc.cached(realtime_data_new_user_cnt, endTime, userCnt)
}

// TODO 暂无
func (svc *RealtimeSvc) cachedOnlineUserCnt() (err error) {

	return nil
}

// 00:00:00 数据放到 前一天
func (svc *RealtimeSvc) cached(dateType string, endTime time.Time, value any) (err error) {
	endTime = endTime.Add(time.Millisecond)
	rKey := realtimeDataRKeyFormatByTime(dateType, endTime)
	var endTimeStr = endTime.Format(pkg.TIME_FORMAT)
	if endTimeStr == "00:00:00" {
		endTimeStr = "24:00:00" // 改为24:00:00 方便排序
		rKey = realtimeDataRKeyFormatByTime(dateType, endTime.AddDate(0, 0, -1))
	}

	err = svc.rdb.HSet(svc.ctx, rKey, endTimeStr, value).Err()
	if err != nil {
		svc.logger.Errorf("%s, rdb.HSet: %v", dateType, err)
		return err
	}
	err = svc.rdb.Expire(svc.ctx, rKey, 72*time.Hour).Err()
	if err != nil {
		svc.logger.Errorf("%s, rdb.Expire: %v", dateType, err)
		return err
	}

	return nil
}

func (svc *RealtimeSvc) All() (tData, yData map[string][][2]string, summaryData map[string]interface{}, e *errcode.Error) {
	cTime := time.Now()

	eg := errgroup.Group{}
	eg.Go(func() (err error) {
		tData, err = svc.all(cTime)
		return err
	})

	eg.Go(func() (err error) {
		yData, err = svc.all(cTime.AddDate(0, 0, -1))
		return err
	})

	eg.Go(func() (err error) {
		summaryData, err = svc.listSummary(cTime)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	return tData, yData, summaryData, nil
}

func (svc *RealtimeSvc) all(cTime time.Time) (data map[string][][2]string, err error) {
	data = make(map[string][][2]string)
	dataChan := make(chan [2]any)

	eg := errgroup.Group{}
	eg.Go(func() error { return svc.allActive(cTime, dataChan) })
	eg.Go(func() error { return svc.allPating(cTime, dataChan) })
	eg.Go(func() error { return svc.allPay(cTime, dataChan) })
	eg.Go(func() error { return svc.allRecharge(cTime, dataChan) })
	eg.Go(func() error { return svc.allDraw(cTime, dataChan) })
	eg.Go(func() error { return svc.allNewUserCnt(cTime, dataChan) })
	eg.Go(func() error { return svc.allOnlineUserCnt(cTime, dataChan) })

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for keyValue := range dataChan {
			data[keyValue[0].(string)] = keyValue[1].([][2]string)
		}

		wg.Done()
	}()

	err = eg.Wait()
	close(dataChan)
	wg.Wait()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (svc *RealtimeSvc) allActive(cTime time.Time, dataChan chan [2]any) (err error) {
	rKey := realtimeDataRKeyFormatByTime(realtime_data_active, cTime)
	data, err := svc.rdb.HGetAll(svc.ctx, rKey).Result()
	if err != nil {
		svc.logger.Errorf("allActive: %v", err)
		return err
	}

	dataChan <- [2]any{realtime_data_active, format(data)}

	return nil
}

func (svc *RealtimeSvc) allPating(cTime time.Time, dataChan chan [2]any) (err error) {
	dataTypeList := []string{"", pating_type_first_prize, pating_type_gashapon, pating_type_chao, pating_type_hole, pating_type_market, pating_type_deliver, pating_type_shop}
	for _, dataType := range dataTypeList {
		rKey := realtimeDataRKeyFormatByTime(realtime_data_pating+dataType, cTime)
		data, err := svc.rdb.HGetAll(svc.ctx, rKey).Result()
		if err != nil {
			svc.logger.Errorf("allPating, %s: %v", dataType, err)
			return err
		}

		dataChan <- [2]any{realtime_data_pating + dataType, format(data)}
	}

	return nil
}

func (svc *RealtimeSvc) allPay(cTime time.Time, dataChan chan [2]any) (err error) {
	for _, dataType := range []string{realtime_data_pay_user_cnt, realtime_data_pay_amount} {
		rKey := realtimeDataRKeyFormatByTime(dataType, cTime)
		data, err := svc.rdb.HGetAll(svc.ctx, rKey).Result()
		if err != nil {
			svc.logger.Errorf("allPay, %s: %v", dataType, err)
			return err
		}

		result := format(data)

		var current decimal.Decimal
		var summaryList [][2]string
		for _, item := range result {
			current = util.Add2Decimal(current, item[1])
			summaryList = append(summaryList, [2]string{item[0], current.String()})
		}

		dataChan <- [2]any{dataType, result}
		dataChan <- [2]any{dataType + realtime_data_summary, summaryList}
	}

	return nil
}

func (svc *RealtimeSvc) allRecharge(cTime time.Time, dataChan chan [2]any) (err error) {
	for _, dataType := range []string{realtime_data_recharge_user_cnt, realtime_data_recharge_amount} {
		rKey := realtimeDataRKeyFormatByTime(dataType, cTime)
		data, err := svc.rdb.HGetAll(svc.ctx, rKey).Result()
		if err != nil {
			svc.logger.Errorf("allRecharge, %s: %v", dataType, err)
			return err
		}

		result := format(data)

		var current decimal.Decimal
		var summaryList [][2]string
		for _, item := range result {
			current = util.Add2Decimal(current, item[1])
			summaryList = append(summaryList, [2]string{item[0], current.String()})
		}

		dataChan <- [2]any{dataType, result}
		dataChan <- [2]any{dataType + realtime_data_summary, summaryList}
	}

	return nil
}

func (svc *RealtimeSvc) allDraw(cTime time.Time, dataChan chan [2]any) (err error) {
	for _, dataType := range []string{realtime_data_draw_user_cnt, realtime_data_draw_amount} {
		rKey := realtimeDataRKeyFormatByTime(dataType, cTime)
		data, err := svc.rdb.HGetAll(svc.ctx, rKey).Result()
		if err != nil {
			svc.logger.Errorf("allDraw, %s: %v", dataType, err)
			return err
		}

		result := format(data)

		var current decimal.Decimal
		var summaryList [][2]string
		for _, item := range result {
			current = util.Add2Decimal(current, item[1])
			summaryList = append(summaryList, [2]string{item[0], current.String()})
		}

		dataChan <- [2]any{dataType, result}
		dataChan <- [2]any{dataType + realtime_data_summary, summaryList}
	}

	return nil
}

func (svc *RealtimeSvc) allNewUserCnt(cTime time.Time, dataChan chan [2]any) (err error) {
	rKey := realtimeDataRKeyFormatByTime(realtime_data_new_user_cnt, cTime)
	data, err := svc.rdb.HGetAll(svc.ctx, rKey).Result()
	if err != nil {
		svc.logger.Errorf("allNewUserCnt: %v", err)
		return err
	}

	result := format(data)

	var current decimal.Decimal
	var summaryList [][2]string
	for _, item := range result {
		current = util.Add2Decimal(current, item[1])
		summaryList = append(summaryList, [2]string{item[0], current.String()})
	}

	dataChan <- [2]any{realtime_data_new_user_cnt, result}
	dataChan <- [2]any{realtime_data_new_user_cnt + realtime_data_summary, summaryList}

	return nil
}

func (svc *RealtimeSvc) allOnlineUserCnt(cTime time.Time, dataChan chan [2]any) (err error) {
	rKey := realtimeDataRKeyFormatByTime(realtime_data_online_user_cnt, cTime)
	data, err := svc.rdb.HGetAll(svc.ctx, rKey).Result()
	if err != nil {
		svc.logger.Errorf("allOnlineUserCnt: %v", err)
		return err
	}

	dataChan <- [2]any{realtime_data_online_user_cnt, format(data)}

	return nil
}

func format(data map[string]string) (result [][2]string) {
	for key, value := range data {
		result = append(result, [2]string{key, value})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i][0] < result[j][0]
	})

	return result
}

func (svc *RealtimeSvc) listSummary(cTime time.Time) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	newUserCntSummary, err := svc.dao.GetNewUserCnt(time.Time{}, cTime)
	if err != nil {
		return nil, err
	}
	payDataSummary, err := svc.dao.GetPayData(time.Time{}, cTime)
	if err != nil {
		return nil, err
	}
	rechargeDataSummary, err := svc.dao.GetRechargeData(time.Time{}, cTime)
	if err != nil {
		return nil, err
	}
	drawDataSummary, err := svc.dao.GetDrawData(time.Time{}, cTime)
	if err != nil {
		return nil, err
	}

	cDate := time.Date(cTime.Year(), cTime.Month(), cTime.Day(), 0, 0, 0, 0, pkg.Location)
	newUserCntToday, err := svc.dao.GetNewUserCnt(cDate, cTime)
	if err != nil {
		return nil, err
	}
	payDataToday, err := svc.dao.GetPayData(cDate, cTime)
	if err != nil {
		return nil, err
	}
	rechargeDataToday, err := svc.dao.GetRechargeData(cDate, cTime)
	if err != nil {
		return nil, err
	}

	data["newUserCntSummary"] = newUserCntSummary
	data["payUserCntSummary"] = payDataSummary["user_cnt"]
	data["payAmountSummary"] = util.ConvertAmount2Decimal(payDataSummary["amount"])
	data["rechargeUserCntSummary"] = rechargeDataSummary["user_cnt"]
	data["rechargeAmountSummary"] = util.ConvertAmount2Decimal(rechargeDataSummary["amount"])
	data["drawAmountSummary"] = util.ConvertAmount2Decimal(drawDataSummary["amount"])
	data["payAmountARPPU"] = util.SaveDivide2Decimal(data["payAmountSummary"], payDataSummary["user_cnt"]).Round(2)
	data["newUserCntSummaryToday"] = newUserCntToday
	data["payAmountSummaryToday"] = util.ConvertAmount2Decimal(payDataToday["amount"])
	data["rechargeAmountSummaryToday"] = util.ConvertAmount2Decimal(rechargeDataToday["amount"])

	return data, nil
}
