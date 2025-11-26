package api

import (
	"context"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/balance/form"
	"data_backend/apps/v2/internal/inquire/balance/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type BalanceApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewBalanceApi() *BalanceApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".BalanceApi")))
	return &BalanceApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (api *BalanceApi) OptionsSourceType(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewBalanceSvc(ctx, local.CMSDB, local.CenterDB, local.RedisClient, api.logger, local.NewAlarm)
	data := svc.OptionsSourceType()

	response.ToResponseData(data)
}

func (api *BalanceApi) OptionsChannelType(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewBalanceSvc(ctx, local.CMSDB, local.CenterDB, local.RedisClient, api.logger, local.NewAlarm)
	data := svc.OptionsChannelType()

	response.ToResponseData(data)
}

func (api *BalanceApi) OptionsPaySourceType(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewBalanceSvc(ctx, local.CMSDB, local.CenterDB, local.RedisClient, api.logger, local.NewAlarm)
	data := svc.OptionsPaySourceType()

	response.ToResponseData(data)
}

func (api *BalanceApi) List(ctx *gin.Context) {
	params := &form.ListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewBalanceSvc(ctx, local.CMSDB, local.CenterDB, local.RedisClient, api.logger, local.NewAlarm)
	summary, data, err := svc.List(params)
	if err != nil {
		api.logger.Errorf("List: %v", err)
		response.ToErrorResponse(err)
		return
	}
	response.ToResponse(gin.H{
		"data": data,
		"headers": map[string]any{
			"summary": summary,
		},
	})
}

func (api *BalanceApi) AddComment(ctx *gin.Context) {
	params := &form.AddCommentRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewBalanceSvc(ctx, local.CMSDB, local.CenterDB, local.RedisClient, api.logger, local.NewAlarm)
	err := svc.AddComment(params)
	if err != nil {
		api.logger.Errorf("AddComment: %v", err)
		response.ToErrorResponse(err)
		return
	}
	response.ToResponseOK()
}

func (api *BalanceApi) DeleteComment(ctx *gin.Context) {
	params := &form.DeleteCommentRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewBalanceSvc(ctx, local.CMSDB, local.CenterDB, local.RedisClient, api.logger, local.NewAlarm)
	err := svc.DeleteComment(params)
	if err != nil {
		api.logger.Errorf("DeleteComment: %v", err)
		response.ToErrorResponse(err)
		return
	}
	response.ToResponseOK()
}

func (api *BalanceApi) Export(ctx *gin.Context) {
	params := &form.AllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewBalanceSvc(ctx, local.CMSDB, local.CenterDB, local.RedisClient, api.logger, local.NewAlarm)
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

func (api *BalanceApi) OptionsBalanceType(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewBalanceSvc(ctx, local.CMSDB, local.CenterDB, local.RedisClient, api.logger, local.NewAlarm)
	data := svc.OptionsBalanceType()

	response.ToResponseData(data)
}
