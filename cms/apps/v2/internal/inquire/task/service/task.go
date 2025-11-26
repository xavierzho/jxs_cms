package service

import (
	"context"
	"data_backend/apps/v2/internal/inquire/task/dao"
	"data_backend/apps/v2/internal/inquire/task/form"
	"data_backend/internal/global"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	dao    *dao.TaskDao
}

func NewTaskSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *TaskSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".TaskSvc")))
	return &TaskSvc{
		ctx:    ctx,
		logger: log,
		dao:    dao.NewTaskDao(center, log),
	}
}

func (svc *TaskSvc) OptionsType() []map[string]string {
	return []map[string]string{
		{"value": "1", "label": global.I18n.T(svc.ctx.Request.Context(), "task.type", "1")},
		{"value": "2", "label": global.I18n.T(svc.ctx.Request.Context(), "task.type", "2")},
		{"value": "3", "label": global.I18n.T(svc.ctx.Request.Context(), "task.type", "3")},
		{"value": "4", "label": global.I18n.T(svc.ctx.Request.Context(), "task.type", "4")},
		{"value": "5", "label": global.I18n.T(svc.ctx.Request.Context(), "task.type", "5")},
	}
}

func (svc *TaskSvc) OptionsKey() []map[string]string {
	return []map[string]string{
		{"value": "CostAmount", "label": global.I18n.T(svc.ctx.Request.Context(), "task.key", "CostAmount")},
		{"value": "PrizeValue", "label": global.I18n.T(svc.ctx.Request.Context(), "task.key", "PrizeValue")},
		{"value": "Week1", "label": global.I18n.T(svc.ctx.Request.Context(), "task.key", "Week1")},
		{"value": "Week2", "label": global.I18n.T(svc.ctx.Request.Context(), "task.key", "Week2")},
		{"value": "Week3", "label": global.I18n.T(svc.ctx.Request.Context(), "task.key", "Week3")},
		{"value": "Week4", "label": global.I18n.T(svc.ctx.Request.Context(), "task.key", "Week4")},
		{"value": "Week5", "label": global.I18n.T(svc.ctx.Request.Context(), "task.key", "Week5")},
		{"value": "Week6", "label": global.I18n.T(svc.ctx.Request.Context(), "task.key", "Week6")},
		{"value": "Week7", "label": global.I18n.T(svc.ctx.Request.Context(), "task.key", "Week7")},
		{"value": "Weekend", "label": global.I18n.T(svc.ctx.Request.Context(), "task.key", "Weekend")},
	}
}

// 任务记录列表
func (svc *TaskSvc) List(params *form.ListRequest) (summary map[string]any, data []*form.TaskSvc, e *errcode.Error) {
	dateTimeRange, paramsGroup, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}
	summary, _data, err := svc.dao.GetList(dateTimeRange, paramsGroup, &params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	summary, data, err = form.Format(svc.ctx.Request.Context(), summary, _data)
	if err != nil {
		return nil, nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}

// 奖励详情
func (svc *TaskSvc) GetAwardDetail(params *form.DetailRequest) (data any, e *errcode.Error) {

	_data, err := svc.dao.GetAwardDetail(params.TaskID)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.FormatAwardDetail(svc.ctx.Request.Context(), _data)
	if err != nil {
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}
	return
}

func (svc *TaskSvc) ExportList(params *form.ListRequest) (data *excel.Excel[*form.TaskSvc], e *errcode.Error) {
	dateTimeRange, paramsGroup, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_, _data, err := svc.dao.GetList(dateTimeRange, paramsGroup, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.FormatList2Excel(svc.ctx.Request.Context(), dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("ExportLog FormatLog2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}

func (svc *TaskSvc) ExportAwardDetail(params *form.ListRequest) (data *excel.Excel[*form.TaskAwardDetailSvc], e *errcode.Error) {
	dateTimeRange, paramsGroup, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.GetListAwardDetail(dateTimeRange, paramsGroup, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.FormatTaskAwardDetail2Excel(svc.ctx.Request.Context(), dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("ExportLog FormatTaskAwardDetail2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
