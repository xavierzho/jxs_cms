package service

import (
	"context"
	"errors"
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

type RecallSvc struct {
	logger  *logger.Logger
	dao     *dao.RecallDao
	userDao *cDao.UserDao
}

func NewRecallSvc(ctx *gin.Context, engine, center *gorm.DB, log *logger.Logger) *RecallSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".RecallSvc")))
	return &RecallSvc{
		logger:  log,
		dao:     dao.NewRecallDao(engine, center, log),
		userDao: cDao.NewUserDao(center, log), //注入user dao访问对象
	}
}

func (svc *RecallSvc) Generate(params *form.GenerateRequest) (e *errcode.Error) {
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

func (svc *RecallSvc) generate(cDate time.Time, queryParams database.QueryWhereGroup) (err error) {
	data, err := svc._generate(cDate, queryParams)
	if err != nil {
		return err
	}

	if err = svc.dao.Save(data...); err != nil {
		return err
	}

	return
}

func (svc *RecallSvc) _generate(cDate time.Time, queryParams database.QueryWhereGroup) (data []*dao.Recall, err error) {
	data, err = svc.dao.Generate(cDate, queryParams)
	if err != nil {
		return nil, err
	}

	return
}

// 获取user信息
func (svc *RecallSvc) userInfo(params *form.AllRequest, queryParams database.QueryWhereGroup) (listUserCondition, generateUserCondition database.QueryWhereGroup, err error) {
	if params.UserType == form.UserType_None || len(queryParams) == 0 {
		return nil, nil, nil
	}

	user, err := svc.userDao.First(queryParams)
	if err != nil {
		return nil, nil, err
	}

	if params.UserType == form.UserType_User {
		listUserCondition = database.QueryWhereGroup{{Prefix: "user_id = ?", Value: []any{user.UserID}}}
		generateUserCondition = database.QueryWhereGroup{{Prefix: "ui.user_id = ?", Value: []any{user.UserID}}}
	} else if params.UserType == form.UserType_ParentUser {
		listUserCondition = database.QueryWhereGroup{{Prefix: "parent_user_id = ?", Value: []any{user.UserID}}}
		generateUserCondition = database.QueryWhereGroup{{Prefix: "ui.parent_user_id = ?", Value: []any{user.UserID}}}
	}

	return
}

func (svc *RecallSvc) List(params *form.ListRequest) (summary map[string]any, data []*form.Recall, e *errcode.Error) {
	dateRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pkg.Location)
	dateRangeCondition := database.QueryWhereGroup{
		{Prefix: "date <> ?", Value: []any{today.Format(pkg.DATE_FORMAT)}},
	}

	listUserCondition, generateUserCondition, err := svc.userInfo(&params.AllRequest, queryParams)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil
	}

	condition := append(dateRangeCondition, listUserCondition...)

	summaryField := []string{
		"sum(amount) as amount",
		"sum(point) as point",
		"count(0) as total",
	}

	// 排除掉当天的数据 若 需要包含当天数据再加回到 summary, 避免出错
	_summary, _data, err := svc.dao.ListAndSummary(summaryField, dateRange, condition, params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	if !dateRange[0].After(today) && !dateRange[1].Before(today) {
		todayData, err := svc._generate(today, generateUserCondition)
		if err != nil {
			return nil, nil, errcode.QueryFail.WithDetails(err.Error())
		}
		_data = append(todayData, _data...)
		var amount = decimal.Zero
		var point = decimal.Zero
		for _, item := range todayData {
			amount = amount.Add(convert.GetDecimal(item.Amount))
			point = amount.Add(convert.GetDecimal(item.Point))
		}
		_summary["amount"] = convert.GetDecimal(_summary["amount"]).Add(amount)
		_summary["point"] = convert.GetDecimal(_summary["point"]).Add(point)
		_summary["total"] = convert.GetInt(_summary["total"]) + len(todayData)

	}

	summary, data = form.Format(dateRange, _summary, _data)

	return
}

func (svc *RecallSvc) Export(params *form.AllRequest) (data *excel.Excel[*form.Recall], e *errcode.Error) {
	dateRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pkg.Location)
	dateRangeCondition := database.QueryWhereGroup{
		{Prefix: "date <> ?", Value: []any{today.Format(pkg.DATE_FORMAT)}},
	}

	listUserCondition, generateUserCondition, err := svc.userInfo(params, queryParams)
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
		todayData, err := svc._generate(today, generateUserCondition)
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
