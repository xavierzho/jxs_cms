package service

import (
	"context"
	"time"

	"data_backend/apps/v2/internal/report/market/dao"
	"data_backend/apps/v2/internal/report/market/form"
	iDao "data_backend/internal/dao"
	iErrcode "data_backend/internal/errcode"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MarketSvc struct {
	logger *logger.Logger
	dao    *dao.MarketDao
}

func NewMarketSvc(ctx *gin.Context, engine, center *gorm.DB, log *logger.Logger) *MarketSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".MarketSvc")))
	return &MarketSvc{
		logger: log,
		dao:    dao.NewMarketDao(engine, center, log),
	}
}

func (svc *MarketSvc) Generate(params *form.GenerateRequest) (e *errcode.Error) {
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

func (svc *MarketSvc) generate(cDate time.Time) (err error) {
	data, err := svc._generate(cDate)
	if err != nil {
		return err
	}

	if err = svc.dao.Save(data); err != nil {
		return err
	}

	return
}

func (svc *MarketSvc) _generate(cDate time.Time) (data *dao.Market, err error) {
	dataCnt, dataAmount, err := svc.dao.Generate(cDate)
	if err != nil {
		return nil, err
	}

	data = &dao.Market{
		DailyModel: iDao.DailyModel{
			Date: cDate.Format(pkg.DATE_FORMAT),
		},
		UserCnt:  dataCnt.UserCnt,
		OrderCnt: dataCnt.OrderCnt,
		Amount0:  dataAmount.Amount0,
		Amount1:  dataAmount.Amount1,
	}

	return
}

func (svc *MarketSvc) List(params *form.ListRequest) (summary map[string]any, data []form.Market, e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pkg.Location)

	summaryFiled := []string{
		"sum(order_cnt) as order_cnt",
		"sum(amount_0) as amount_0",
		"sum(amount_1) as amount_1",
	}
	// 排除掉当天的数据 若 需要包含当天数据再加回到 summary, 避免出错
	_summary, _data, err := svc.dao.ListAndSummary(summaryFiled, dateRange, database.QueryWhereGroup{{Prefix: "date <> ?", Value: []any{today.Format(pkg.DATE_FORMAT)}}}, params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if !dateRange[0].After(today) && !dateRange[1].Before(today) {
		todayData, err := svc._generate(today)
		if err != nil {
			return nil, nil, errcode.QueryFail.WithDetails(err.Error())
		}
		_data = append(_data, todayData)

		_summary["order_cnt"] = convert.GetDecimal(_summary["order_cnt"]).Add(convert.GetDecimal(todayData.OrderCnt))
		_summary["amount_0"] = convert.GetDecimal(_summary["amount_0"]).Add(convert.GetDecimal(todayData.Amount0))
		_summary["amount_1"] = convert.GetDecimal(_summary["amount_1"]).Add(convert.GetDecimal(todayData.Amount1))
	}

	// 将时间进行分页
	dateRange = params.Pager.PaginationDateRange(dateRange)
	_summary["total"] = params.Pager.TotalRows

	summary, data, err = form.Format(dateRange, _summary, _data)
	if err != nil {
		svc.logger.Errorf("List, Format: %v", err)
		return nil, nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}

func (svc *MarketSvc) Export(params *form.AllRequest) (data *excel.Excel[form.Market], e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pkg.Location)

	_data, err := svc.dao.All(dateRange, database.QueryWhereGroup{{Prefix: "date <> ?", Value: []any{today.Format(pkg.DATE_FORMAT)}}})
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if !dateRange[0].After(today) && !dateRange[1].Before(today) {
		todayData, err := svc._generate(today)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		_data = append(_data, todayData)
	}

	data, err = form.Format2Excel(dateRange, _data)
	if err != nil {
		svc.logger.Errorf("Export Format2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
