package service

import (
	"context"
	"time"

	"data_backend/apps/v2/internal/report/dashboard/dao"
	"data_backend/apps/v2/internal/report/dashboard/form"
	marketDao "data_backend/apps/v2/internal/report/market/dao"
	"data_backend/internal/app"
	iDao "data_backend/internal/dao"
	iErrcode "data_backend/internal/errcode"
	"data_backend/pkg"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardSvc struct {
	logger    *logger.Logger
	dao       *dao.DashboardDao
	marketDao *marketDao.MarketDao
}

func NewDashboardSvc(ctx *gin.Context, engine, center *gorm.DB, log *logger.Logger) *DashboardSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".DashboardSvc")))
	return &DashboardSvc{
		logger:    log,
		dao:       dao.NewDashboardDao(engine, center, log),
		marketDao: marketDao.NewMarketDao(engine, center, log),
	}
}

func (svc *DashboardSvc) Generate(params *form.GenerateRequest) (e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return errcode.InvalidParams.WithDetails(err.Error())
	}

	for cDate := dateRange[0]; !cDate.After(dateRange[1]); cDate = cDate.AddDate(0, 0, 1) {
		if err = svc.generate(cDate); err != nil {
			return iErrcode.SQLExecFail.WithDetails(err.Error())
		}
	}

	return nil
}

func (svc *DashboardSvc) generate(cDate time.Time) (err error) {
	data, err := svc._generate(cDate, cDate.Add(24*time.Hour-time.Millisecond))
	if err != nil {
		return err
	}
	if err = svc.dao.Save(data); err != nil {
		return err
	}

	return nil
}

func (svc *DashboardSvc) _generate(startTime, endTime time.Time) (data *dao.Dashboard, err error) {
	dataGroup, err := svc.dao.Generate(startTime, endTime)
	if err != nil {
		return nil, err
	}

	return &dao.Dashboard{
		DailyModel: iDao.DailyModel{
			Date: endTime.Format(pkg.DATE_FORMAT),
		},
		NewUserCnt:           dataGroup.NewUser.NewUserCnt,
		ActiveUserCnt:        dataGroup.ActiveUser.ActiveUserCnt,
		PatingUserCnt:        dataGroup.Pating.PatingUserCnt,
		PatingUserCntNew:     dataGroup.Pating.PatingUserCntNew,
		PayUserCnt:           dataGroup.Pay.PayUserCnt,
		PayUserCntNew:        dataGroup.Pay.PayUserCntNew,
		RechargeUserCnt:      dataGroup.Recharge.RechargeUserCnt,
		RechargeUserCntNew:   dataGroup.Recharge.RechargeUserCntNew,
		RechargeAmount:       dataGroup.Recharge.RechargeAmount,
		RechargeAmountWeChat: dataGroup.Recharge.RechargeAmountWeChat,
		RechargeAmountAli:    dataGroup.Recharge.RechargeAmountAli,
		DrawAmount:           dataGroup.Draw.DrawAmount + dataGroup.RechargeRefund.DrawAmount + dataGroup.SavingRefund.DrawAmount,
	}, nil
}

func (svc *DashboardSvc) List() (data []*form.Dashboard, dataSummary *form.Dashboard, e *errcode.Error) {
	cTime := time.Now()
	currentDate := time.Date(cTime.Year(), cTime.Month(), cTime.Day(), 0, 0, 0, 0, pkg.Location)
	startDate := currentDate.AddDate(0, 0, -9)

	// history
	dataHistory, _, err := svc.dao.List([2]time.Time{startDate, currentDate.AddDate(0, 0, -1)}, nil, app.Pager{Page: 1, PageSize: 10})
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}
	// today
	dataToday, err := svc._generate(currentDate, currentDate.Add(24*time.Hour-time.Millisecond))
	if err != nil {
		return nil, nil, iErrcode.SQLExecFail.WithDetails(err.Error())
	}

	// merge market data
	marketData, err := svc.marketDao.All([2]time.Time{startDate, currentDate.Add(24*time.Hour - time.Millisecond)}, nil)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}
	// market today
	marketCntToday, marketAmountToday, err := svc.marketDao.Generate(currentDate)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}
	marketData = append(marketData, &marketDao.Market{
		DailyModel: iDao.DailyModel{Date: currentDate.Format(pkg.DATE_FORMAT)},
		OrderCnt:   marketCntToday.OrderCnt,
		Amount0:    marketAmountToday.Amount0,
		Amount1:    marketAmountToday.Amount1,
	})

	dashboardList := append(dataHistory, dataToday)
	for _, item := range dashboardList {
		for _, mItem := range marketData {
			if item.Date == mItem.Date {
				item.MarketOrderCnt = int(mItem.OrderCnt)
				item.MarketAmount0 = mItem.Amount0
				item.MarketAmount1 = mItem.Amount1
				break
			}
		}

	}

	data = form.Format([2]time.Time{startDate, currentDate}, dashboardList...)

	// summary
	_dataSummary, err := svc._generate(time.Time{}, currentDate.Add(24*time.Hour-time.Millisecond))
	if err != nil {
		return nil, nil, iErrcode.SQLExecFail.WithDetails(err.Error())
	}
	// summary market
	if _dataSummary != nil {
		for _, mItem := range marketData {
			_dataSummary.MarketOrderCnt += int(mItem.OrderCnt)
			_dataSummary.MarketAmount0 += mItem.Amount0
			_dataSummary.MarketAmount1 += mItem.Amount1
		}
	}

	dataSummaryList := form.Format([2]time.Time{currentDate, currentDate}, _dataSummary)
	if len(dataSummaryList) == 0 {
		return nil, nil, errcode.TransformFail.WithDetails("dataSummaryList length = 0")
	}
	dataSummary = dataSummaryList[0]

	return
}
