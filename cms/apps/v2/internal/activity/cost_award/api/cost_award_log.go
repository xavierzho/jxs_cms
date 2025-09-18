package api

import (
	"context"

	"data_backend/apps/v2/internal/activity/cost_award/form"
	"data_backend/apps/v2/internal/activity/cost_award/service"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type CostAwardLogApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewCostAwardLogApi() *CostAwardLogApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".CostAwardLogApi")))
	return &CostAwardLogApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}

}

func (api *CostAwardLogApi) OptionsLogType(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewCostAwardLogSvc(ctx, local.CenterDB, api.logger)
	data := svc.OptionsLogType()

	response.ToResponseData(data)
}

func (api *CostAwardLogApi) List(ctx *gin.Context) {
	params := &form.ListLogRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewCostAwardLogSvc(ctx, local.CenterDB, api.logger)
	summary, data, err := svc.List(params)
	if err != nil {
		api.logger.Errorf("List: %v", err)
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

func (api *CostAwardLogApi) Export(ctx *gin.Context) {
	params := &form.AllLogRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewCostAwardLogSvc(ctx, local.CenterDB, api.logger)
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
