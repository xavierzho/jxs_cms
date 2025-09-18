package api

import (
	"context"

	"data_backend/apps/v2/internal/admin/service"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/app"
	iService "data_backend/internal/service"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type MenuApi struct {
	logger *logger.Logger
}

func NewMenuApi() MenuApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".MenuApi")))
	return MenuApi{
		logger: log,
	}
}

func (api *MenuApi) All(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := iService.NewMenuSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	data, err := svc.All(service.MenuList())
	if err != nil {
		api.logger.Errorf("All: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}
