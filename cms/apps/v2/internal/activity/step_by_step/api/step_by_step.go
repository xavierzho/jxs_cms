package api

import (
	"context"

	"data_backend/apps/v2/internal/activity/step_by_step/form"
	"data_backend/apps/v2/internal/activity/step_by_step/service"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type StepByStepApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewStepByStepApi() *StepByStepApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".StepByStepApi")))
	return &StepByStepApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}

}

func (api *StepByStepApi) LogList(ctx *gin.Context) {
	params := &form.LogListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewStepByStepSvc(ctx, local.CenterDB, api.logger)
	data, summary, err := svc.LogList(params)
	if err != nil {
		api.logger.Errorf("LogList: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponse(gin.H{
		"data": data,
		"headers": map[string]any{
			"summary": summary,
			"total":   summary["total"],
		},
	})
}

func (api *StepByStepApi) LogExport(ctx *gin.Context) {
	params := &form.LogAllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewStepByStepSvc(ctx, local.CenterDB, api.logger)
	excelModel, err := svc.LogExport(params)
	if err != nil {
		api.logger.Errorf("LogExport: %v", err)
		response.ToErrorResponse(err)
		return
	}

	e := response.ExportFile(ctx, excelModel.Excelize, excelModel.FileName)
	if e != nil {
		api.logger.Errorf("response.ExportFile err: %v", e.Error())
	}
}

func (api *StepByStepApi) Detail(ctx *gin.Context) {
	params := &form.DetailRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewStepByStepSvc(ctx, local.CenterDB, api.logger)
	data, err := svc.Detail(params)
	if err != nil {
		api.logger.Errorf("Detail: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}

func (api *StepByStepApi) DetailExport(ctx *gin.Context) {
	params := &form.DetailAllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewStepByStepSvc(ctx, local.CenterDB, api.logger)
	excelModel, err := svc.DetailExport(params)
	if err != nil {
		api.logger.Errorf("DetailExport: %v", err)
		response.ToErrorResponse(err)
		return
	}

	e := response.ExportFile(ctx, excelModel.Excelize, excelModel.FileName)
	if e != nil {
		api.logger.Errorf("response.ExportFile err: %v", e.Error())
	}
}
