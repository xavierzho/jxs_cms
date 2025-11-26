package api

import (
	"context"

	"data_backend/apps/v2/internal/activity/redemption_code/form"
	"data_backend/apps/v2/internal/activity/redemption_code/service"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type RedemptionCodeApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewRedemptionCodeApi() *RedemptionCodeApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".RedemptionCodeApi")))
	return &RedemptionCodeApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (api *RedemptionCodeApi) Log(ctx *gin.Context) {
	params := &form.LogRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewRedemptionCodeSvc(ctx, local.CenterDB, api.logger)
	summary, data, err := svc.Log(params)
	if err != nil {
		api.logger.Errorf("Log: %v", err)
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

func (api *RedemptionCodeApi) GetAwardDetail(ctx *gin.Context) {
	params := &form.DetailRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewRedemptionCodeSvc(ctx, local.CenterDB, api.logger)
	data, err := svc.GetAwardDetail(params)
	if err != nil {
		api.logger.Errorf("GetDetail: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}

func (api *RedemptionCodeApi) ExportLog(ctx *gin.Context) {
	params := &form.LogRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewRedemptionCodeSvc(ctx, local.CenterDB, api.logger)
	excelModel, err := svc.ExportList(params)
	if err != nil {
		api.logger.Errorf("ExportLog: %v", err)
		response.ToErrorResponse(err)
		return
	}

	e := response.ExportFile(ctx, excelModel.Excelize, excelModel.FileName)
	if e != nil {
		api.logger.Errorf("response.ExportFile err: %v", e.Error())
	}
}

func (api *RedemptionCodeApi) ExportAwardDetail(ctx *gin.Context) {
	params := &form.LogRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewRedemptionCodeSvc(ctx, local.CenterDB, api.logger)
	excelModel, err := svc.ExportAwardDetail(params)
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
