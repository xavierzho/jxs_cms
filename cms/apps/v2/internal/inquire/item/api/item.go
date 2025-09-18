package api

import (
	"context"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/item/form"
	"data_backend/apps/v2/internal/inquire/item/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type ItemApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewItemApi() *ItemApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".ItemApi")))
	return &ItemApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (api *ItemApi) OptionsLogType(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewItemSvc(ctx, local.CenterDB, api.logger)
	data := svc.OptionsLogType()

	response.ToResponseData(data)
}

func (api *ItemApi) GetLog(ctx *gin.Context) {
	params := &form.LogRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewItemSvc(ctx, local.CenterDB, api.logger)
	summary, data, err := svc.GetLog(params)
	if err != nil {
		api.logger.Errorf("GetLog: %v", err)
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

func (api *ItemApi) ExportLog(ctx *gin.Context) {
	params := &form.LogAllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewItemSvc(ctx, local.CenterDB, api.logger)
	excelModel, err := svc.ExportLog(params)
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

func (api *ItemApi) GetDetail(ctx *gin.Context) {
	params := &form.DetailRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewItemSvc(ctx, local.CenterDB, api.logger)
	data, err := svc.GetDetail(params)
	if err != nil {
		api.logger.Errorf("GetDetail: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}

func (api *ItemApi) ExportDetail(ctx *gin.Context) {
	params := &form.DetailAllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewItemSvc(ctx, local.CenterDB, api.logger)
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
