package api

import (
	"context"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/app"
	iService "data_backend/internal/service"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type PermissionApi struct {
	logger *logger.Logger
}

func NewPermissionApi() PermissionApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".PermissionApi")))
	return PermissionApi{
		logger: log,
	}
}

func (api *PermissionApi) Options(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := iService.NewPermissionSvc(ctx, local.CMSDB, api.logger, local.NewAlarm)
	data, err := svc.Options()
	if err != nil {
		api.logger.Errorf("Options: %v", err)
		response.ToErrorResponse(errcode.QueryFail.WithDetails(err.Error()))
		return
	}

	response.ToResponseData(data)
}
