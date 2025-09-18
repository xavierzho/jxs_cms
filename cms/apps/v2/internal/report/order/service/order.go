package service

import (
	"context"
	"errors"
	"time"

	cDao "data_backend/apps/v2/internal/common/dao"
	"data_backend/apps/v2/internal/report/order/dao"
	"data_backend/apps/v2/internal/report/order/form"
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

type DeliveryOrderSvc struct {
	logger  *logger.Logger
	dao     *dao.DeliveryOrderDao
	userDao *cDao.UserDao
}

func NewDeliveryOrderSvc(ctx *gin.Context, engine, center *gorm.DB, log *logger.Logger) *DeliveryOrderSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".DeliveryOrderSvc")))
	return &DeliveryOrderSvc{
		logger:  log,
		dao:     dao.NewDeliveryOrderDao(engine, center, log),
		userDao: cDao.NewUserDao(center, log), //注入user dao访问对象
	}
}

func (svc *DeliveryOrderSvc) Generate(params *form.GenerateRequest) (e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return errcode.InvalidParams.WithDetails(err.Error())
	}

	for cDate := dateRange[0]; !cDate.After(dateRange[1]); cDate = cDate.AddDate(0, 0, 1) {
		if err = svc.generate(cDate, nil); err != nil {
			return iErrcode.SQLExecFail.WithDetails(err.Error())
		}
	}

	return nil
}

func (svc *DeliveryOrderSvc) generate(cDate time.Time, queryParams database.QueryWhereGroup) (err error) {
	data, err := svc._generate(cDate, queryParams)
	if err != nil {
		return err
	}

	if err = svc.dao.Save(data...); err != nil {
		return err
	}

	return
}

func (svc *DeliveryOrderSvc) _generate(cDate time.Time, queryParams database.QueryWhereGroup) (data []*dao.DeliveryOrder, err error) {
	data, err = svc.dao.Generate(cDate, queryParams)
	if err != nil {
		return nil, err
	}

	return
}

// 获取user信息
func (svc *DeliveryOrderSvc) userInfo(queryParams database.QueryWhereGroup) (listUserCondition database.QueryWhereGroup, err error) {
	if len(queryParams) == 0 {
		return nil, nil
	}
	user, err := svc.userDao.First(queryParams)
	if err != nil {
		return nil, err
	}

	listUserCondition = database.QueryWhereGroup{{Prefix: "user_id = ?", Value: []any{user.UserID}}}

	return
}

func (svc *DeliveryOrderSvc) List(params *form.ListRequest) (summary map[string]any, data []*form.DeliveryOrder, e *errcode.Error) {
	dateRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pkg.Location)
	dateRangeCondition := database.QueryWhereGroup{
		{Prefix: "date <> ?", Value: []any{today.Format(pkg.DATE_FORMAT)}},
	}

	listUserCondition, err := svc.userInfo(queryParams)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil
	}

	condition := append(dateRangeCondition, listUserCondition...)

	summaryField := []string{
		"sum(show_price) as show_price",
		"sum(inner_price) as inner_price",
		"sum(recycling_price) as recycling_price",
		"count(0) as total",
	}

	// 排除掉当天的数据 若 需要包含当天数据再加回到 summary, 避免出错
	_summary, _data, err := svc.dao.ListAndSummary(summaryField, dateRange, condition, params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}
	if !dateRange[0].After(today) && !dateRange[1].Before(today) {
		todayData, err := svc._generate(today, listUserCondition)
		if err != nil {
			return nil, nil, errcode.QueryFail.WithDetails(err.Error())
		}
		_data = append(todayData, _data...)
		var showPrice, innerPrice, recyclingPrice = decimal.Zero, decimal.Zero, decimal.Zero
		for _, item := range todayData {
			showPrice = showPrice.Add(convert.GetDecimal(item.ShowPrice))
			innerPrice = innerPrice.Add(convert.GetDecimal(item.InnerPrice))
			recyclingPrice = recyclingPrice.Add(convert.GetDecimal(item.RecyclingPrice))
		}
		_summary["show_price"] = convert.GetDecimal(_summary["show_price"]).Add(showPrice)
		_summary["inner_price"] = convert.GetDecimal(_summary["inner_price"]).Add(innerPrice)
		_summary["recycling_price"] = convert.GetDecimal(_summary["recycling_price"]).Add(recyclingPrice)
		_summary["total"] = convert.GetInt(_summary["total"]) + len(todayData)

	}

	summary, data = form.Format(dateRange, _summary, _data)

	return
}

func (svc *DeliveryOrderSvc) Export(params *form.AllRequest) (data *excel.Excel[*form.DeliveryOrder], e *errcode.Error) {
	dateRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pkg.Location)
	dateRangeCondition := database.QueryWhereGroup{
		{Prefix: "date <> ?", Value: []any{today.Format(pkg.DATE_FORMAT)}},
	}

	listUserCondition, err := svc.userInfo(queryParams)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	condition := append(dateRangeCondition, listUserCondition...)

	_data, err := svc.dao.All(dateRange, condition)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if !dateRange[0].After(today) && !dateRange[1].Before(today) {
		todayData, err := svc._generate(today, listUserCondition)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		_data = append(todayData, _data...)
	}

	data, err = form.Format2Excel(dateRange, _data)
	if err != nil {
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
