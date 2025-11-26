package api

import (
	"context"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/task/form"
	"data_backend/apps/v2/internal/inquire/task/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type TaskApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewTaskApi() *TaskApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".TaskApi")))
	return &TaskApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (api *TaskApi) OptionsType(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewTaskSvc(ctx, local.CenterDB, api.logger)
	data := svc.OptionsType()

	response.ToResponseData(data)
}

func (api *TaskApi) OptionsKey(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewTaskSvc(ctx, local.CenterDB, api.logger)
	data := svc.OptionsKey()

	response.ToResponseData(data)
}

func (api *TaskApi) List(ctx *gin.Context) {
	params := &form.ListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewTaskSvc(ctx, local.CenterDB, api.logger)
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

func (api *TaskApi) GetAwardDetail(ctx *gin.Context) {
	params := &form.DetailRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewTaskSvc(ctx, local.CenterDB, api.logger)
	data, err := svc.GetAwardDetail(params)
	if err != nil {
		api.logger.Errorf("GetDetail: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}

func (api *TaskApi) ExportList(ctx *gin.Context) {
	params := &form.ListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewTaskSvc(ctx, local.CenterDB, api.logger)
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

func (api *TaskApi) ExportAwardDetail(ctx *gin.Context) {
	params := &form.ListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewTaskSvc(ctx, local.CenterDB, api.logger)
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
