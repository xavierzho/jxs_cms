package api

import (
	"context"
	"fmt"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/bet/form"
	"data_backend/apps/v2/internal/report/bet/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type BetApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewBetApi() *BetApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".BetApi")))
	return &BetApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (api *BetApi) Generate(ctx *gin.Context) {
	params := &form.GenerateRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewBetSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	go func() {
		if err := svc.Generate(params); err != nil {
			api.alarm.AlertErrorMsg(fmt.Sprintf("BetSvc.Generate: %v", err), message.CMS_ID)
		}
	}()

	response.ToResponseOK()
}

func (api *BetApi) All(ctx *gin.Context) {
	params := &form.AllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewBetSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	data, err := svc.All(params)
	if err != nil {
		api.logger.Errorf("All: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}
