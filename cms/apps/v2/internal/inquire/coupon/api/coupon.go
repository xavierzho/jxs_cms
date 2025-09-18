package api

import (
	"context"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/coupon/form"
	"data_backend/apps/v2/internal/inquire/coupon/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type CouponApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewCouponApi() *CouponApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".CouponApi")))
	return &CouponApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (api *CouponApi) OptionsCouponType(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewCouponSvc(ctx, local.CenterDB, api.logger)
	data := svc.OptionsCouponType()

	response.ToResponseData(data)
}

func (api *CouponApi) OptionsCouponActionType(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewCouponSvc(ctx, local.CenterDB, api.logger)
	data := svc.OptionsCouponActionType()

	response.ToResponseData(data)
}

func (api *CouponApi) List(ctx *gin.Context) {
	params := &form.ListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewCouponSvc(ctx, local.CenterDB, api.logger)
	data, count, err := svc.List(params)
	if err != nil {
		api.logger.Errorf("List: %v", err)
		response.ToErrorResponse(err)
		return
	}
	response.ToResponse(gin.H{
		"data": data,
		"headers": map[string]any{
			"total": count,
		},
	})
}

func (api *CouponApi) Export(ctx *gin.Context) {
	params := &form.AllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewCouponSvc(ctx, local.CenterDB, api.logger)
	excelModel, err := svc.Export(params)
	if err != nil {
		api.logger.Errorf("Export: %v", err)
		response.ToErrorResponse(err)
		return
	}

	e := response.ExportFile(ctx, excelModel.Excelize, excelModel.FileName)
	if e != nil {
		api.logger.Errorf("response.ExportFile err: %v", e.Error())
	}
}
