package service

import (
	"context"
	"fmt"
	"time"

	"data_backend/apps/v2/internal/report/revenue/dao"
	"data_backend/apps/v2/internal/report/revenue/form"
	iDao "data_backend/internal/dao"
	iErrcode "data_backend/internal/errcode"
	"data_backend/pkg"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type RevenueSvc struct {
	ctx    *gin.Context
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewRevenueSvc(ctx *gin.Context, engine, center *gorm.DB, log *logger.Logger) *RevenueSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".RevenueSvc")))
	return &RevenueSvc{
		ctx:    ctx,
		engine: engine,
		center: center,
		logger: log,
	}
}

func (svc *RevenueSvc) Generate(params *form.GenerateRequest) (e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return errcode.InvalidParams.WithDetails(err.Error())
	}

	for cDate := dateRange[0]; !cDate.After(dateRange[1]); cDate = cDate.AddDate(0, 0, 1) {
		if e = svc.generate(params.DataTypeList, cDate); e != nil {
			return e
		}
	}

	return nil
}

func (svc *RevenueSvc) generate(dataTypeList []string, cDate time.Time) (e *errcode.Error) {
	eg := errgroup.Group{}
	for _, dataType := range dataTypeList {
		switch dataType {
		case form.REVENUE_DATA_TYPE_PAY:
			eg.Go(func() error { return svc.generatePay(cDate) })
		case form.REVENUE_DATA_TYPE_DRAW:
			eg.Go(func() error { return svc.generateDraw(cDate) })
		case form.REVENUE_DATA_TYPE_BALANCE:
			eg.Go(func() error { return svc.generateBalance(cDate) })
		case form.REVENUE_DATA_TYPE_ACTIVE:
			eg.Go(func() error { return svc.generateActive(cDate) })
		case form.REVENUE_DATA_TYPE_PATING:
			eg.Go(func() error { return svc.generatePating(cDate) })
		case form.REVENUE_DATA_TYPE_WASTAGE:
			eg.Go(func() error { return svc.generateWastage(cDate) })
		default:
			_dataType := dataType
			eg.Go(func() error { return fmt.Errorf("not expected data_type: " + _dataType) })
		}
	}

	if err := eg.Wait(); err != nil {
		return iErrcode.SQLExecFail.WithDetails(err.Error())
	}

	return nil
}

func (svc *RevenueSvc) generatePay(cDate time.Time) (err error) {
	payDao := dao.NewPayDao(svc.engine, svc.center, svc.logger)
	data, err := svc._generatePay(payDao, cDate)
	if err != nil {
		return err
	}

	if err = payDao.Save(data); err != nil {
		return err
	}

	return nil
}

// 集市, 发货, 存在退款(潮币)情况。统计消费金额时 按完成时间统计 所以这里不再统计 退款(潮币)
func (svc *RevenueSvc) _generatePay(payDao *dao.PayDao, cDate time.Time) (data *dao.Pay, err error) {
	payGroup, err := payDao.Generate(cDate)
	if err != nil {
		return nil, err
	}

	return &dao.Pay{
		DailyModel: iDao.DailyModel{
			Date: cDate.Format(pkg.DATE_FORMAT),
		},
		Amount:       payGroup.Pay.Amount,
		AmountBet:    payGroup.Pay.AmountBet,
		AmountNew:    payGroup.Pay.AmountNew,
		AmountOld:    payGroup.Pay.AmountOld,
		UserCnt:      payGroup.Pay.UserCnt,
		UserCntNew:   payGroup.Pay.UserCntNew,
		UserCntOld:   payGroup.Pay.UserCntOld,
		UserCntFirst: payGroup.Pay.UserCntFirst,
		// RefundAmount:         payGroup.Refund.RefundAmount,
		// RefundUserCnt:        payGroup.Refund.RefundUserCnt,
		RechargeAmount:             payGroup.Recharge.RechargeAmount,
		RechargeAmountWeChat:       payGroup.Recharge.RechargeAmountWeChat,
		RechargeAmountAli:          payGroup.Recharge.RechargeAmountAli,
		RechargeRefundAmount:       payGroup.RechargeRefund.RechargeRefundAmount,
		RechargeRefundAmountWeChat: payGroup.RechargeRefund.RechargeRefundAmountWeChat,
		RechargeRefundAmountAli:    payGroup.RechargeRefund.RechargeRefundAmountAli,
		DiscountAmount:             payGroup.Recharge.DiscountAmount,
		DiscountAmountWeChat:       payGroup.Recharge.DiscountAmountWeChat,
		DiscountAmountAli:          payGroup.Recharge.DiscountAmountAli,
		SavingAmount:               payGroup.Saving.SavingAmount,
		SavingAmountWeChat:         payGroup.Saving.SavingAmountWeChat,
		SavingAmountAli:            payGroup.Saving.SavingAmountAli,
		SavingRefundAmount:         payGroup.SavingRefund.SavingRefundAmount,
		SavingRefundAmountWeChat:   payGroup.SavingRefund.SavingRefundAmountWeChat,
		SavingRefundAmountAli:      payGroup.SavingRefund.SavingRefundAmountAli,
	}, nil
}

func (svc *RevenueSvc) generateDraw(cDate time.Time) (err error) {
	drawDao := dao.NewDrawDao(svc.engine, svc.center, svc.logger)
	data, err := svc._generateDraw(drawDao, cDate)
	if err != nil {
		return err
	}

	if err = drawDao.Save(data); err != nil {
		return err
	}

	return nil
}

func (svc *RevenueSvc) _generateDraw(drawDao *dao.DrawDao, cDate time.Time) (data *dao.Draw, err error) {
	data, err = drawDao.Generate(cDate)
	if err != nil {
		return nil, err
	}

	data.Date = cDate.Format(pkg.DATE_FORMAT)

	return data, nil
}

func (svc *RevenueSvc) generateBalance(cDate time.Time) (err error) {
	balanceDao := dao.NewBalanceDao(svc.engine, svc.center, svc.logger)
	data, err := svc._generateBalance(balanceDao, cDate)
	if err != nil {
		return err
	}

	if err = balanceDao.Save(data); err != nil {
		return err
	}

	return nil
}

func (svc *RevenueSvc) _generateBalance(balanceDao *dao.BalanceDao, cDate time.Time) (data *dao.Balance, err error) {
	data, err = balanceDao.Generate()
	if err != nil {
		return nil, err
	}

	data.Date = cDate.Format(pkg.DATE_FORMAT)

	return data, nil
}

func (svc *RevenueSvc) generateActive(cDate time.Time) (err error) {
	activeDao := dao.NewActiveDao(svc.engine, svc.center, svc.logger)
	data, err := svc._generateActive(activeDao, cDate)
	if err != nil {
		return err
	}

	if err = activeDao.Save(data); err != nil {
		return err
	}

	return nil
}

func (svc *RevenueSvc) _generateActive(activeDao *dao.ActiveDao, cDate time.Time) (data *dao.Active, err error) {
	dataRegister, dataActive, dataValidated, err := activeDao.Generate(cDate)
	if err != nil {
		return nil, err
	}

	return &dao.Active{
		DailyModel: iDao.DailyModel{
			Date: cDate.Format(pkg.DATE_FORMAT),
		},
		ActivateCnt:    dataActive.ActivateCnt,
		ActivateCntNew: dataActive.ActivateCntNew,
		RegisterCnt:    dataRegister.RegisterCnt,
		ActiveCnt:      dataActive.ActiveCnt,
		ActiveCntNew:   dataActive.ActiveCntNew,
		ActiveCntOld:   dataActive.ActiveCntOld,
		MaxOnlineCnt:   0,
		ValidatedCnt7:  dataValidated.ValidatedCnt7,
		ValidatedCnt15: dataValidated.ValidatedCnt15,
	}, nil
}

func (svc *RevenueSvc) generatePating(cDate time.Time) (err error) {
	patingDao := dao.NewPatingDao(svc.engine, svc.center, svc.logger)
	data, err := svc._generatePating(patingDao, cDate)
	if err != nil {
		return err
	}

	if err = patingDao.Save(data); err != nil {
		return err
	}

	return nil
}

func (svc *RevenueSvc) _generatePating(patingDao *dao.PatingDao, cDate time.Time) (data *dao.Pating, err error) {
	data, err = patingDao.Generate(cDate)
	if err != nil {
		return nil, err
	}

	data.Date = cDate.Format(pkg.DATE_FORMAT)

	return data, nil
}

func (svc *RevenueSvc) generateWastage(cDate time.Time) (err error) {
	targetTime := time.Now().AddDate(0, 0, -6)
	if cDate.After(targetTime) { // 开始时间未到
		return nil
	}

	wastageDao := dao.NewWastageDao(svc.engine, svc.center, svc.logger)
	data, err := wastageDao.Generate(cDate)
	if err != nil {
		return err
	}

	data.Date = cDate.Format(pkg.DATE_FORMAT)

	if err = wastageDao.Save(data); err != nil {
		return err
	}

	return nil
}

func (svc *RevenueSvc) All(params *form.AllRequest) (data interface{}, e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pkg.Location)
	hasToday := !dateRange[0].After(today) && !dateRange[1].Before(today)
	switch params.DataType {
	case form.REVENUE_REPORT_TYPE_ACTIVE:
		return svc.allActive(dateRange, hasToday, today)
	case form.REVENUE_REPORT_TYPE_PATING:
		return svc.allPating(dateRange, hasToday, today)
	case form.REVENUE_REPORT_TYPE_PAY:
		return svc.allPay(dateRange, hasToday, today)
	case form.REVENUE_REPORT_TYPE_DRAW:
		return svc.allDraw(dateRange, hasToday, today)
	case form.REVENUE_REPORT_TYPE_SUMMARY:
		return svc.allSummary(dateRange, hasToday, today)
	default:
		return nil, errcode.InvalidParams.WithDetails("not expected data_type: " + params.DataType)
	}
}

func (svc *RevenueSvc) allActive(dateRange [2]time.Time, hasToday bool, today time.Time) (data []form.Active, e *errcode.Error) {
	activeDao := dao.NewActiveDao(svc.engine, svc.center, svc.logger)
	activeData, err := activeDao.All(dateRange, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	wastageDao := dao.NewWastageDao(svc.engine, svc.center, svc.logger)
	wastageData, err := wastageDao.All(dateRange, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if hasToday {
		activeTodayData, err := svc._generateActive(activeDao, today)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		activeData = append(activeData, activeTodayData) // 填充到最后用于替换缓存中的
	}

	data, err = form.FormatActive(dateRange, activeData, wastageData)
	if err != nil {
		svc.logger.Errorf("allActive, FormatActive: %v", err)
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return data, nil
}

func (svc *RevenueSvc) allPating(dateRange [2]time.Time, hasToday bool, today time.Time) (data []form.Pating, e *errcode.Error) {
	patingDao := dao.NewPatingDao(svc.engine, svc.center, svc.logger)
	patingData, err := patingDao.All(dateRange, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	activeDao := dao.NewActiveDao(svc.engine, svc.center, svc.logger)
	activeData, err := activeDao.All(dateRange, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if hasToday {
		patingTodayData, err := svc._generatePating(patingDao, today)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		patingData = append(patingData, patingTodayData)

		activeTodayData, err := svc._generateActive(activeDao, today)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		activeData = append(activeData, activeTodayData) // 填充到最后用于替换缓存中的
	}

	data, err = form.FormatPating(dateRange, patingData, activeData)
	if err != nil {
		svc.logger.Errorf("allPating, FormatPating: %v", err)
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return data, nil
}

func (svc *RevenueSvc) allPay(dateRange [2]time.Time, hasToday bool, today time.Time) (data []form.Pay, e *errcode.Error) {
	payDao := dao.NewPayDao(svc.engine, svc.center, svc.logger)
	payData, err := payDao.All(dateRange, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	activeDao := dao.NewActiveDao(svc.engine, svc.center, svc.logger)
	activeData, err := activeDao.All(dateRange, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if hasToday {
		payTodayData, err := svc._generatePay(payDao, today)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		payData = append(payData, payTodayData)

		activeTodayData, err := svc._generateActive(activeDao, today)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		activeData = append(activeData, activeTodayData) // 填充到最后用于替换缓存中的
	}

	data, err = form.FormatPay(dateRange, payData, activeData)
	if err != nil {
		svc.logger.Errorf("allPay, FormatPay: %v", err)
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return data, nil
}

func (svc *RevenueSvc) allDraw(dateRange [2]time.Time, hasToday bool, today time.Time) (data []form.Draw, e *errcode.Error) {
	drawDao := dao.NewDrawDao(svc.engine, svc.center, svc.logger)
	drawData, err := drawDao.All(dateRange, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	activeDao := dao.NewActiveDao(svc.engine, svc.center, svc.logger)
	activeData, err := activeDao.All(dateRange, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if hasToday {
		drawTodayData, err := svc._generateDraw(drawDao, today)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		drawData = append(drawData, drawTodayData)

		activeTodayData, err := svc._generateActive(activeDao, today)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		activeData = append(activeData, activeTodayData) // 填充到最后用于替换缓存中的
	}

	data, err = form.FormatDraw(dateRange, drawData, activeData)
	if err != nil {
		svc.logger.Errorf("allDraw, FormatDraw: %v", err)
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return data, nil
}

func (svc *RevenueSvc) allSummary(dateRange [2]time.Time, hasToday bool, today time.Time) (data []form.Summary, e *errcode.Error) {
	summaryDao := dao.NewRevenueDao(svc.engine, svc.logger)
	summaryData, err := summaryDao.All(dateRange)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if hasToday {
		var activeTodayData *dao.Active
		var payTodayData *dao.Pay
		var drawTodayData *dao.Draw
		var balanceTodayData *dao.Balance
		eg := errgroup.Group{}
		eg.Go(func() (err error) {
			activeDao := dao.NewActiveDao(svc.engine, svc.center, svc.logger)
			activeTodayData, err = svc._generateActive(activeDao, today)
			return err
		})

		eg.Go(func() (err error) {
			payDao := dao.NewPayDao(svc.engine, svc.center, svc.logger)
			payTodayData, err = svc._generatePay(payDao, today)
			return err
		})

		eg.Go(func() (err error) {
			drawDao := dao.NewDrawDao(svc.engine, svc.center, svc.logger)
			drawTodayData, err = svc._generateDraw(drawDao, today)
			return err
		})

		eg.Go(func() (err error) {
			balanceDao := dao.NewBalanceDao(svc.engine, svc.center, svc.logger)
			balanceTodayData, err = svc._generateBalance(balanceDao, today)
			return err
		})

		err = eg.Wait()
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}

		var todayStr = today.Format(pkg.DATE_FORMAT)
		var summaryTodayData = map[string]any{
			"date":                          todayStr,
			"wallet_balance":                balanceTodayData.Wallet,
			"merchant_balance":              balanceTodayData.Merchant,
			"gold_balance":                  balanceTodayData.Gold,
			"pay_amount":                    payTodayData.Amount,
			"pay_amount_bet":                payTodayData.AmountBet,
			"recharge_amount":               payTodayData.RechargeAmount,
			"recharge_amount_wechat":        payTodayData.RechargeAmountWeChat,
			"recharge_amount_ali":           payTodayData.RechargeAmountAli,
			"recharge_refund_amount":        payTodayData.RechargeRefundAmount,
			"recharge_refund_amount_wechat": payTodayData.RechargeRefundAmountWeChat,
			"recharge_refund_amount_ali":    payTodayData.RechargeRefundAmountAli,
			"draw_amount":                   drawTodayData.Amount,
			"tax_amount":                    0,
			"refund_amount":                 payTodayData.RefundAmount,
			"active_cnt":                    activeTodayData.ActiveCnt,
			"discount_amount":               payTodayData.DiscountAmount,
			"discount_amount_wechat":        payTodayData.DiscountAmountWeChat,
			"discount_amount_ali":           payTodayData.DiscountAmountAli,
			"saving_amount":                 payTodayData.SavingAmount,
			"saving_amount_wechat":          payTodayData.SavingAmountWeChat,
			"saving_amount_ali":             payTodayData.SavingAmountAli,
			"saving_refund_amount":          payTodayData.SavingRefundAmount,
			"saving_refund_amount_wechat":   payTodayData.SavingRefundAmountWeChat,
			"saving_refund_amount_ali":      payTodayData.SavingRefundAmountAli,
		}

		summaryData = append(summaryData, summaryTodayData)
	}

	data, err = form.FormatSummary(dateRange, summaryData)
	if err != nil {
		svc.logger.Errorf("allSummary, FormatSummary: %v", err)
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return data, nil
}
