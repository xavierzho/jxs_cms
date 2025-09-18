package api

import (
	"context"
	"fmt"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/dashboard/form"
	"data_backend/apps/v2/internal/report/dashboard/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type DashboardApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewDashboardApi() *DashboardApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".DashboardApi")))
	return &DashboardApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}

}

func (api *DashboardApi) Generate(ctx *gin.Context) {
	params := &form.GenerateRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewDashboardSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	go func() {
		if err := svc.Generate(params); err != nil {
			api.alarm.AlertErrorMsg(fmt.Sprintf("DashboardSvc.Generate: %v", err), message.CMS_ID)
		}
	}()

	response.ToResponseOK()
}

func (api *DashboardApi) List(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := service.NewDashboardSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	data, summary, err := svc.List()
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
