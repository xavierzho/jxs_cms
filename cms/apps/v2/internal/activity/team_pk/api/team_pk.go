package api

import (
	"context"

	"data_backend/apps/v2/internal/activity/team_pk/form"
	"data_backend/apps/v2/internal/activity/team_pk/service"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type TeamPKApi struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewTeamPKApi() *TeamPKApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".TeamPKApi")))
	return &TeamPKApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}

}

func (api *TeamPKApi) List(ctx *gin.Context) {
	params := &form.ListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewTeamPKSvc(ctx, local.CenterDB, api.logger)
	data, total, err := svc.List(params)
	if err != nil {
		api.logger.Errorf("List: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponse(gin.H{
		"data": data,
		"headers": map[string]any{
			"total": total,
		},
	})
}

func (api *TeamPKApi) Export(ctx *gin.Context) {
	params := &form.AllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewTeamPKSvc(ctx, local.CenterDB, api.logger)
	excelModel, err := svc.Export(params)
	if err != nil {
		api.logger.Errorf("Export: %v", err)
		response.ToErrorResponse(err)
		return
	}

	e := response.ExportFile(ctx, excelModel.Excelize, excelModel.FileName)
	if e != nil {
		api.logger.Errorf("response.ExportFile err: %v", e.Error())
	}
}
