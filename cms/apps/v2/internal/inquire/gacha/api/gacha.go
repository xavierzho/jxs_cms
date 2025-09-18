package api

import (
	"context"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/gacha/form"
	"data_backend/apps/v2/internal/inquire/gacha/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type GachaApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewGachaApi() *GachaApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".GachaApi")))
	return &GachaApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (api *GachaApi) OptionsGachaType(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewGachaSvc(ctx, local.CenterDB, api.logger)
	data := svc.OptionsGachaType()

	response.ToResponseData(data)
}

func (api *GachaApi) GetRevenue(ctx *gin.Context) {
	params := &form.RevenueRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewGachaSvc(ctx, local.CenterDB, api.logger)
	summary, data, err := svc.GetRevenue(params)
	if err != nil {
		api.logger.Errorf("GetRevenue: %v", err)
		response.ToErrorResponse(err)
		return
	}
	response.ToResponse(gin.H{
		"data": data,
		"headers": map[string]any{
			"total":   summary["total"],
			"summary": summary,
		},
	})
}

func (api *GachaApi) GetDetail(ctx *gin.Context) {
	params := &form.DetailRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewGachaSvc(ctx, local.CenterDB, api.logger)
	data, err := svc.GetDetail(params)
	if err != nil {
		api.logger.Errorf("GetDetail: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}

func (api *GachaApi) ExportDetail(ctx *gin.Context) {
	params := &form.DetailRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewGachaSvc(ctx, local.CenterDB, api.logger)
	excelModel, err := svc.ExportDetail(params)
	if err != nil {
		api.logger.Errorf("ExportDetail: %v", err)
		response.ToErrorResponse(err)
		return
	}

	e := response.ExportFile(ctx, excelModel.Excelize, excelModel.FileName)
	if e != nil {
		api.logger.Errorf("response.ExportFile err: %v", e.Error())
	}
}
