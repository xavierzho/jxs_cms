package service

import (
	"context"
	"time"

	"data_backend/apps/v2/internal/report/dashboard/dao"
	"data_backend/apps/v2/internal/report/dashboard/form"
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
	logger *logger.Logger
	dao    *dao.DashboardDao
}

func NewDashboardSvc(ctx *gin.Context, engine, center *gorm.DB, log *logger.Logger) *DashboardSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".DashboardSvc")))
	return &DashboardSvc{
		logger: log,
		dao:    dao.NewDashboardDao(engine, center, log),
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
		PatingUserCnt:        dataGroup.PatingData.PatingUserCnt,
		PatingUserCntNew:     dataGroup.PatingData.PatingUserCntNew,
		PayUserCnt:           dataGroup.PayData.PayUserCnt,
		PayUserCntNew:        dataGroup.PayData.PayUserCntNew,
		RechargeUserCnt:      dataGroup.RechargeData.RechargeUserCnt,
		RechargeUserCntNew:   dataGroup.RechargeData.RechargeUserCntNew,
		RechargeAmount:       dataGroup.RechargeData.RechargeAmount,
		RechargeAmountWeChat: dataGroup.RechargeData.RechargeAmountWeChat,
		RechargeAmountAli:    dataGroup.RechargeData.RechargeAmountAli,
		DrawAmount:           dataGroup.DrawData.DrawAmount,
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

	data = form.Format([2]time.Time{startDate, currentDate}, append(dataHistory, dataToday)...)

	// summary
	_dataSummary, err := svc._generate(time.Time{}, currentDate.Add(24*time.Hour-time.Millisecond))
	if err != nil {
		return nil, nil, iErrcode.SQLExecFail.WithDetails(err.Error())
	}
	dataSummaryList := form.Format([2]time.Time{currentDate, currentDate}, _dataSummary)
	if len(dataSummaryList) == 0 {
		return nil, nil, errcode.TransformFail.WithDetails("dataSummaryList length = 0")
	}
	dataSummary = dataSummaryList[0]

	return
}
