package service

import (
	"context"

	"data_backend/apps/v2/internal/inquire/gacha/dao"
	"data_backend/apps/v2/internal/inquire/gacha/form"
	"data_backend/internal/global"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GachaSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	dao    *dao.GachaDao
}

func NewGachaSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *GachaSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".GachaSvc")))
	return &GachaSvc{
		ctx:    ctx,
		logger: log,
		dao:    dao.NewGachaDao(center, log),
	}
}

func (svc *GachaSvc) OptionsGachaType() []map[string]string {
	return []map[string]string{
		{"value": "101", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "101")},
		{"value": "102", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "102")},
		{"value": "103", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "103")},
		{"value": "104", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "104")},
	}
}

func (svc *GachaSvc) GetRevenue(params *form.RevenueRequest) (summary map[string]any, data []*form.GachaRevenue, e *errcode.Error) {
	paramsGroup, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	summary, _data, err := svc.dao.GetRevenue(paramsGroup, params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	summary, data = form.FormatRevenue(svc.ctx.Request.Context(), summary, _data)

	return
}

func (svc *GachaSvc) GetDetail(params *form.DetailRequest) (data []*form.GachaDetail, e *errcode.Error) {
	queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.GetDetail(queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data = form.FormatDetail(_data)

	return
}

func (svc *GachaSvc) ExportDetail(params *form.DetailRequest) (data *excel.Excel[*form.GachaDetail], e *errcode.Error) {
	queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.GetDetail(queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.Format2Excel(params, _data)
	if err != nil {
		svc.logger.Errorf("ExportDetail Format2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
