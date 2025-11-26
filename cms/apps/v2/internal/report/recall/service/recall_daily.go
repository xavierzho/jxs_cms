package service

import (
	"context"
	"time"

	cDao "data_backend/apps/v2/internal/common/dao"
	"data_backend/apps/v2/internal/report/recall/dao"
	"data_backend/apps/v2/internal/report/recall/form"
	iErrcode "data_backend/internal/errcode"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RecallDailySvc struct {
	logger  *logger.Logger
	dao     *dao.RecallDailyDao
	userDao *cDao.UserDao
}

func NewRecallDailySvc(ctx *gin.Context, engine, center *gorm.DB, log *logger.Logger) *RecallDailySvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".RecallDailySvc")))
	return &RecallDailySvc{
		logger:  log,
		dao:     dao.NewRecallDailyDao(engine, center, log),
		userDao: cDao.NewUserDao(center, log),
	}
}

func (svc *RecallDailySvc) Generate(params *form.GenerateDailyRequest) (e *errcode.Error) {
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

func (svc *RecallDailySvc) generate(cDate time.Time) (err error) {
	data, err := svc.dao.Generate(cDate)
	if err != nil {
		return err
	}

	if err = svc.dao.Save(data...); err != nil {
		return err
	}

	return
}

func (svc *RecallDailySvc) List(params *form.ListDailyRequest) (summary map[string]any, data []*form.RecallDaily, e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pkg.Location)

	summaryField := []string{
		"sum(total_amount) as total_amount",
		"sum(amount) as amount",
		"sum(difference) as difference",
		"count(0) as total",
	}
	_summary, _data, err := svc.dao.ListAndSummary(summaryField, dateRange, database.QueryWhereGroup{{Prefix: "date <> ?", Value: []any{today.Format(pkg.DATE_FORMAT)}}}, params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if !dateRange[0].After(today) && !dateRange[1].Before(today) {
		todayData, err := svc.dao.Generate(today)
		if err != nil {
			return nil, nil, errcode.QueryFail.WithDetails(err.Error())
		}
		_data = append(todayData, _data...)
		var totalAmount, amount, difference = decimal.Zero, decimal.Zero, decimal.Zero
		for _, item := range todayData {
			totalAmount = totalAmount.Add(convert.GetDecimal(item.TotalAmount))
			amount = amount.Add(convert.GetDecimal(item.Amount))
			difference = difference.Add(convert.GetDecimal(item.Difference))
		}
		_summary["total_amount"] = convert.GetDecimal(_summary["total_amount"]).Add(totalAmount)
		_summary["amount"] = convert.GetDecimal(_summary["amount"]).Add(amount)
		_summary["difference"] = convert.GetDecimal(_summary["difference"]).Add(difference)
		_summary["total"] = convert.GetInt(_summary["total"]) + len(todayData)

	}

	summary, data = form.FormatDaily(dateRange, _summary, _data)

	return
}

func (svc *RecallDailySvc) Export(params *form.AllDailyRequest) (data *excel.Excel[*form.RecallDaily], e *errcode.Error) {
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
		todayData, err := svc.dao.Generate(today)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		_data = append(todayData, _data...)
	}

	data, err = form.Format2DailyExcel(dateRange, _data)
	if err != nil {
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
