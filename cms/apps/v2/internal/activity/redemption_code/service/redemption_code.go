package service

import (
	"context"

	"data_backend/apps/v2/internal/activity/redemption_code/dao"
	"data_backend/apps/v2/internal/activity/redemption_code/form"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RedemptionCodeSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	dao    *dao.RedemptionCodeDao
}

func NewRedemptionCodeSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *RedemptionCodeSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".RedemptionCodeSvc")))
	return &RedemptionCodeSvc{
		ctx:    ctx,
		logger: log,
		dao:    dao.NewRedemptionCodeDao(center, log),
	}
}

func (svc *RedemptionCodeSvc) Log(params *form.LogRequest) (summary map[string]any, data []*form.RedemptionCodeSvc, e *errcode.Error) {
	dateTimeRange, paramsGroup, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}
	summary, _data, err := svc.dao.GetLog(dateTimeRange, paramsGroup, &params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	summary, data, err = form.Format(svc.ctx.Request.Context(), summary, _data)
	if err != nil {
		return nil, nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}

func (svc *RedemptionCodeSvc) GetAwardDetail(params *form.DetailRequest) (data any, e *errcode.Error) {

	_data, err := svc.dao.GetAwardDetail(params.LogID)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.FormatAwardDetail(svc.ctx.Request.Context(), _data)
	if err != nil {
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}
	return
}

func (svc *RedemptionCodeSvc) ExportList(params *form.LogRequest) (data *excel.Excel[*form.RedemptionCodeSvc], e *errcode.Error) {
	dateTimeRange, paramsGroup, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_, _data, err := svc.dao.GetLog(dateTimeRange, paramsGroup, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.FormatLog2Excel(svc.ctx.Request.Context(), dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("ExportLog FormatLog2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}

func (svc *RedemptionCodeSvc) ExportAwardDetail(params *form.LogRequest) (data *excel.Excel[*form.RedemptionCodeAwardDetailSvc], e *errcode.Error) {
	dateTimeRange, paramsGroup, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.GetLogAwardDetail(dateTimeRange, paramsGroup, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.FormatAwardDetail2Excel(svc.ctx.Request.Context(), dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("ExportLog FormatTaskAwardDetail2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
