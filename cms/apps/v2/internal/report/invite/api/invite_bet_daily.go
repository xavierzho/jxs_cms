package api

import (
	"context"
	"fmt"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/invite/form"
	"data_backend/apps/v2/internal/report/invite/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type InviteBetDailyApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewInviteBetDailyApi() *InviteBetDailyApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".InviteBetDailyApi")))
	return &InviteBetDailyApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}

}

func (api *InviteBetDailyApi) Generate(ctx *gin.Context) {
	params := &form.GenerateDailyRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewInviteBetDailySvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	go func() {
		if err := svc.Generate(params); err != nil {
			api.alarm.AlertErrorMsg(fmt.Sprintf("InviteBetSvc.Generate: %v", err), message.CMS_ID)
		}
	}()

	response.ToResponseOK()
}

func (api *InviteBetDailyApi) List(ctx *gin.Context) {
	params := &form.ListDailyRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewInviteBetDailySvc(ctx, local.CMSDB, local.CenterDB, api.logger)
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

func (api *InviteBetDailyApi) Export(ctx *gin.Context) {
	params := &form.AllDailyRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewInviteBetDailySvc(ctx, local.CMSDB, local.CenterDB, api.logger)
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
