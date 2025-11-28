package service

import (
	"context"

	"data_backend/apps/v2/internal/activity/cost_award/dao"
	"data_backend/apps/v2/internal/activity/cost_award/form"
	"data_backend/internal/global"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CostAwardLogSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	dao    *dao.CostAwardLogDao
}

func NewCostAwardLogSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *CostAwardLogSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".CostAwardLogSvc")))
	return &CostAwardLogSvc{
		ctx:    ctx,
		logger: log,
		dao:    dao.NewCostAwardLogDao(center, log),
	}
}

func (svc *CostAwardLogSvc) OptionsLogType() []map[string]string {
	return []map[string]string{
		{"value": "100002", "label": global.I18n.T(svc.ctx, "source_type", "100002")},
		{"value": "100008", "label": global.I18n.T(svc.ctx, "source_type", "100008")},
		{"value": "100009", "label": global.I18n.T(svc.ctx, "source_type", "100009")},
		{"value": "100010", "label": global.I18n.T(svc.ctx, "source_type", "100010")},
		{"value": "200000", "label": global.I18n.T(svc.ctx, "source_type", "200000")},
	}
}

func (svc *CostAwardLogSvc) List(params *form.ListLogRequest) (summary map[string]any, data []*form.CostAwardLog, e *errcode.Error) {
	dateTimeRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_summary, _data, err := svc.dao.List(dateTimeRange, queryParams, params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	summary, data = form.FormatLog(svc.ctx, _summary, _data)

	return
}

func (svc *CostAwardLogSvc) Export(params *form.AllLogRequest) (data *excel.Excel[*form.CostAwardLog], e *errcode.Error) {
	dateTimeRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.All(dateTimeRange, queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.FormatLog2Excel(svc.ctx, dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("Export Format2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
