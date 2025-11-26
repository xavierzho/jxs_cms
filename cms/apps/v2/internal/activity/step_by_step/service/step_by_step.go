package service

import (
	"context"

	"data_backend/apps/v2/internal/activity/step_by_step/dao"
	"data_backend/apps/v2/internal/activity/step_by_step/form"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StepByStepSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	dao    *dao.StepByStepDao
}

func NewStepByStepSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *StepByStepSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".StepByStepSvc")))
	return &StepByStepSvc{
		ctx:    ctx,
		logger: log,
		dao:    dao.NewStepByStepDao(center, log),
	}
}

func (svc *StepByStepSvc) LogList(params *form.LogListRequest) (data []*form.StepByStepLog, summary map[string]any, e *errcode.Error) {
	dateTimeRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, _summary, err := svc.dao.LogList(dateTimeRange, queryParams, params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, summary = form.LogFormat(params.PointType, _data, _summary)

	return
}

func (svc *StepByStepSvc) LogExport(params *form.LogAllRequest) (data *excel.Excel[*form.StepByStepLog], e *errcode.Error) {
	dateTimeRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.LogAll(dateTimeRange, queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.LogFormat2Excel(dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("LogExport Format2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}

func (svc *StepByStepSvc) Detail(params *form.DetailRequest) (data []*form.StepByStepAward, e *errcode.Error) {
	queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.Detail(queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.DetailFormat(_data)
	if err != nil {
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}

func (svc *StepByStepSvc) DetailExport(params *form.DetailAllRequest) (data *excel.Excel[*form.StepByStepAward], e *errcode.Error) {
	dateTimeRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.DetailAll(dateTimeRange, queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.DetailFormat2Excel(dateTimeRange, _data)
	if err != nil {
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
