package api

import (
	"context"
	"fmt"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/cohort/form"
	"data_backend/apps/v2/internal/report/cohort/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type CohortApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewCohortApi() *CohortApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".CohortApi")))
	return &CohortApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (api *CohortApi) Generate(ctx *gin.Context) {
	params := &form.GenerateRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewCohortSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	go func() {
		if err := svc.Generate(params); err != nil {
			api.alarm.AlertErrorMsg(fmt.Sprintf("CohortSvc.Generate: %v", err), message.CmsId)
		}
	}()

	response.ToResponseOK()
}

func (api *CohortApi) All(ctx *gin.Context) {
	params := &form.AllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewCohortSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	data, userCnt, err := svc.All(params)
	if err != nil {
		api.logger.Errorf("All: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponse(gin.H{
		"data": data,
		"headers": map[string]interface{}{
			"user_cnt": userCnt,
		},
	})
}
