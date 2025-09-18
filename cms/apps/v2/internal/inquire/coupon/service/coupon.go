package service

import (
	"context"

	"data_backend/apps/v2/internal/inquire/coupon/dao"
	"data_backend/apps/v2/internal/inquire/coupon/form"
	"data_backend/internal/global"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CouponSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	dao    *dao.CouponDao
}

func NewCouponSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *CouponSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".CouponSvc")))
	return &CouponSvc{
		ctx:    ctx,
		logger: log,
		dao:    dao.NewCouponDao(center, log),
	}
}

func (svc *CouponSvc) OptionsCouponType() []map[string]string {
	return []map[string]string{
		{"value": "1", "label": global.I18n.T(svc.ctx.Request.Context(), "coupon.type", "1")},
		{"value": "2", "label": global.I18n.T(svc.ctx.Request.Context(), "coupon.type", "2")},
		{"value": "3", "label": global.I18n.T(svc.ctx.Request.Context(), "coupon.type", "3")},
		{"value": "4", "label": global.I18n.T(svc.ctx.Request.Context(), "coupon.type", "4")},
	}
}

func (svc *CouponSvc) OptionsCouponActionType() []map[string]string {
	return []map[string]string{
		{"value": "0", "label": global.I18n.T(svc.ctx.Request.Context(), "coupon.action", "0")},
		{"value": "1", "label": global.I18n.T(svc.ctx.Request.Context(), "coupon.action", "1")},
		{"value": "2", "label": global.I18n.T(svc.ctx.Request.Context(), "coupon.action", "2")},
	}
}

func (svc *CouponSvc) List(params *form.ListRequest) (data []*form.Coupon, count int64, e *errcode.Error) {
	dateTimeRange, explain, queryParams, err := params.Parse()
	if err != nil {
		return nil, 0, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, count, err := svc.dao.List(dateTimeRange, explain, queryParams, params.Pager)
	if err != nil {
		return nil, 0, errcode.QueryFail.WithDetails(err.Error())
	}

	data = form.Format(svc.ctx, _data)

	return
}

func (svc *CouponSvc) Export(params *form.AllRequest) (data *excel.Excel[*form.Coupon], e *errcode.Error) {
	dateTimeRange, explain, queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.All(dateTimeRange, explain, queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.Format2Excel(svc.ctx, dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("Export Format2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
