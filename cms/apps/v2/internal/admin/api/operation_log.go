package api

import (
	"context"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/app"
	"data_backend/internal/form"
	iService "data_backend/internal/service"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type OperationLogApi struct {
	logger *logger.Logger
}

func NewOperationLogApi() OperationLogApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".OperationLogApi")))
	return OperationLogApi{
		logger: log,
	}
}

func (api *OperationLogApi) List(ctx *gin.Context) {
	params := &form.OperationLogListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := iService.NewOperationLogSvc(ctx, local.CMSDB, api.logger, local.NewAlarm)
	data, count, err := svc.List(params)
	if err != nil {
		api.logger.Errorf("List: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseList(data, count)
}
