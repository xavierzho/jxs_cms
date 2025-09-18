package api

import (
	"context"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/app"
	"data_backend/internal/form"
	iService "data_backend/internal/service"
	"data_backend/pkg/convert"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RoleApi struct {
	logger *logger.Logger
}

func NewRoleApi() RoleApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".RoleApi")))
	return RoleApi{
		logger: log,
	}
}

func (api *RoleApi) Create(ctx *gin.Context) {
	params := &form.RoleCreateRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := iService.NewRoleSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	err := svc.Create(params)
	if err != nil {
		api.logger.Errorf("Create: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseOK()
}

func (api *RoleApi) List(ctx *gin.Context) {
	params := &form.RoleListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := iService.NewRoleSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	data, total, err := svc.List(params)
	if err != nil {
		api.logger.Errorf("List: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseList(data, total)
}

func (api *RoleApi) Update(ctx *gin.Context) {
	params := &form.RoleUpdateRequest{}
	id := convert.StrTo(ctx.Param("id")).MustUInt32()
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := iService.NewRoleSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	data, err := svc.Update(id, params)
	if err != nil {
		api.logger.Errorf("Update: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}

func (api *RoleApi) Options(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := iService.NewRoleSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	data, err := svc.Options()
	if err != nil {
		api.logger.Errorf("Options: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}
