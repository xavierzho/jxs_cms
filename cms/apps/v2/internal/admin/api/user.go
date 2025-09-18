package api

import (
	"context"
	"errors"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/app"
	iErrcode "data_backend/internal/errcode"
	"data_backend/internal/form"
	iService "data_backend/internal/service"
	"data_backend/pkg/convert"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type UserApi struct {
	logger *logger.Logger
}

func NewUserApi() UserApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".UserApi")))
	return UserApi{
		logger: log,
	}
}

func (api *UserApi) Login(ctx *gin.Context) {
	params := &form.LoginRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := iService.NewUserSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	data, err := svc.Login(params)
	if err != nil {
		if !errors.Is(err, iErrcode.UserNotExist) && !errors.Is(err, iErrcode.IncorrectPassword) && !errors.Is(err, iErrcode.UserIsLock) {
			api.logger.Errorf("Login: %v", err)
		}
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}

func (api *UserApi) Create(ctx *gin.Context) {
	params := &form.UserCreateRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := iService.NewUserSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	err := svc.Create(params)
	if err != nil {
		api.logger.Errorf("Create: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseOK()
}

func (api *UserApi) Detail(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := iService.NewUserSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	data, err := svc.CurrentUser()
	if err != nil {
		api.logger.Errorf("CurrentUser: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}

func (api *UserApi) List(ctx *gin.Context) {
	params := &form.UserListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := iService.NewUserSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	data, total, err := svc.List(params)
	if err != nil {
		api.logger.Errorf("List: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseList(data, total)
}

func (api *UserApi) Update(ctx *gin.Context) {
	params := &form.UserUpdateRequest{}
	id := convert.StrTo(ctx.Param("id")).MustUInt32()
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := iService.NewUserSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	err := svc.Update(id, params)
	if err != nil {
		api.logger.Errorf("Update: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseOK()
}

func (api *UserApi) UpdateSelf(ctx *gin.Context) {
	params := &form.UserUpdateSelfRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := iService.NewUserSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	data, err := svc.UpdateSelf(params)
	if err != nil {
		api.logger.Errorf("UpdateSelf: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}

func (api *UserApi) PagePermission(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := iService.NewUserSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	data, err := svc.PagePermission()
	if err != nil {
		api.logger.Errorf("PagePermission: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}

func (api *UserApi) Options(ctx *gin.Context) {
	response := app.NewResponse(ctx)

	svc := iService.NewUserSvc(ctx, local.CMSDB, local.RedisClient, api.logger, local.NewAlarm)
	data, err := svc.Options()
	if err != nil {
		api.logger.Errorf("Options: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponseData(data)
}
