package api

import (
	"context"
	"fmt"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/realtime/form"
	"data_backend/apps/v2/internal/report/realtime/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type RealtimeApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewRealtimeApi() *RealtimeApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".RealtimeApi")))
	return &RealtimeApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}

}

func (api *RealtimeApi) Cached(ctx *gin.Context) {
	params := &form.CachedRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewRealtimeSvc(ctx, local.CenterDB, local.RedisClient, api.logger)
	go func() {
		if err := svc.Cached(params); err != nil {
			api.alarm.AlertErrorMsg(fmt.Sprintf("RealtimeSvc.Cached: %v", err), message.CMS_ID)
		}
	}()

	response.ToResponseOK()
}

func (api *RealtimeApi) All(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewRealtimeSvc(ctx, local.CenterDB, local.RedisClient, api.logger)
	tData, yData, summaryData, err := svc.All()
	if err != nil {
		api.logger.Errorf("All: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponse(gin.H{
		"data": map[string]any{
			"tData":       tData,
			"yData":       yData,
			"summaryData": summaryData,
		},
	})
}
