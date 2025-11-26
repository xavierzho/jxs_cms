package api

import (
	"context"
	"fmt"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/revenue/form"
	"data_backend/apps/v2/internal/report/revenue/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type RevenueApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewRevenueApi() *RevenueApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".RevenueApi")))
	return &RevenueApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}

}

func (api *RevenueApi) Generate(ctx *gin.Context) {
	params := &form.GenerateRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewRevenueSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	go func() {
		if err := svc.Generate(params); err != nil {
			api.alarm.AlertErrorMsg(fmt.Sprintf("RevenueSvc.Generate: %v", err), message.CmsId)
		}
	}()

	response.ToResponseOK()
}

func (api *RevenueApi) All(ctx *gin.Context) {
	params := &form.AllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewRevenueSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	data, err := svc.All(params)
	if err != nil {
		api.logger.Errorf("All: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}
