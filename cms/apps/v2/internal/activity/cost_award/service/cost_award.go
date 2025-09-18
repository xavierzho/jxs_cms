package service

import (
	"context"
	"time"

	"data_backend/apps/v2/internal/activity/cost_award/dao"
	"data_backend/apps/v2/internal/activity/cost_award/form"
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

type CostAwardSvc struct {
	logger *logger.Logger
	dao    *dao.CostAwardDao
}

func NewCostAwardSvc(ctx *gin.Context, engine, center *gorm.DB, log *logger.Logger) *CostAwardSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".CostAwardSvc")))
	return &CostAwardSvc{
		logger: log,
		dao:    dao.NewCostAwardDao(engine, center, log),
	}
}

func (svc *CostAwardSvc) Generate(params *form.GenerateRequest) (e *errcode.Error) {
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

func (svc *CostAwardSvc) generate(cDate time.Time) (err error) {
	data, err := svc._generate(cDate)
	if err != nil {
		return err
	}

	if err = svc.dao.Save(data); err != nil {
		return err
	}

	return
}

func (svc *CostAwardSvc) _generate(cDate time.Time) (data *dao.CostAward, err error) {
	dataLog, dataAward, err := svc.dao.Generate(cDate)
	if err != nil {
		return nil, err
	}

	data = &dao.CostAward{
		DailyModel: iDao.DailyModel{
			Date: cDate.Format(pkg.DATE_FORMAT),
		},
		GetUserCnt:          dataLog.GetUserCnt,
		GetAmount:           dataLog.GetAmount,
		AcceptUserCnt:       dataLog.AcceptUserCnt,
		AcceptAmount:        dataLog.AcceptAmount,
		AwardAmount:         dataAward.AwardAmount,
		AwardItemShowPrice:  dataAward.AwardItemShowPrice,
		AwardItemInnerPrice: dataAward.AwardItemInnerPrice,
	}

	return
}

func (svc *CostAwardSvc) List(params *form.ListRequest) (summary map[string]any, data []*form.CostAward, e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pkg.Location)

	summaryField := []string{
		"sum(get_amount) as get_amount",
		"sum(accept_amount) as accept_amount",
		"sum(award_amount) as award_amount",
		"sum(award_item_show_price) as award_item_show_price",
		"sum(award_item_inner_price) as award_item_inner_price",
	}
	// 排除掉当天的数据 若 需要包含当天数据再加回到 summary, 避免出错
	_summary, _data, err := svc.dao.ListAndSummary(summaryField, dateRange, database.QueryWhereGroup{{Prefix: "date <> ?", Value: []any{today.Format(pkg.DATE_FORMAT)}}}, params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if !dateRange[0].After(today) && !dateRange[1].Before(today) {
		todayData, err := svc._generate(today)
		if err != nil {
			return nil, nil, errcode.QueryFail.WithDetails(err.Error())
		}
		_data = append(_data, todayData)

		_summary["get_amount"] = convert.GetDecimal(_summary["get_amount"]).Add(convert.GetDecimal(todayData.GetAmount))
		_summary["accept_amount"] = convert.GetDecimal(_summary["accept_amount"]).Add(convert.GetDecimal(todayData.AcceptAmount))
		_summary["award_amount"] = convert.GetDecimal(_summary["award_amount"]).Add(convert.GetDecimal(todayData.AwardAmount))
		_summary["award_item_show_price"] = convert.GetDecimal(_summary["award_item_show_price"]).Add(convert.GetDecimal(todayData.AwardItemShowPrice))
		_summary["award_item_inner_price"] = convert.GetDecimal(_summary["award_item_inner_price"]).Add(convert.GetDecimal(todayData.AwardItemInnerPrice))
	}

	// 将时间进行分页
	dateRange = params.Pager.PaginationDateRange(dateRange)
	_summary["total"] = params.Pager.TotalRows

	summary, data = form.Format(dateRange, _summary, _data)

	return
}

func (svc *CostAwardSvc) Export(params *form.AllRequest) (data *excel.Excel[*form.CostAward], e *errcode.Error) {
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
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
