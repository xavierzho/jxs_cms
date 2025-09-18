package api

import (
	"context"

	"data_backend/apps/v2/internal/activity/turntable/form"
	"data_backend/apps/v2/internal/activity/turntable/service"
	"data_backend/apps/v2/internal/common/local"

	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type TurntableApi struct {
	logger *logger.Logger //记录日志
	alarm  message.Alarm  //预警信息
}

// NewTurntableApi 创建并返回一个新的 TurntableApi 实例
func NewTurntableApi() *TurntableApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".TurntableApi")))
	return &TurntableApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}

}

// 转盘抽奖列表
func (api *TurntableApi) List(ctx *gin.Context) {
	params := &form.ListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewTurntableSvc(ctx, local.CenterDB, api.logger)
	summary, data, err := svc.List(params)
	if err != nil {
		api.logger.Errorf("List: %v", err)
		response.ToErrorResponse(err)
		return
	}

	response.ToResponse(gin.H{
		"data": data,
		"headers": map[string]any{
			"total":   summary["total"],
			"summary": summary,
		},
	})
}

func (api *TurntableApi) Export(ctx *gin.Context) {
	params := &form.AllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewTurntableSvc(ctx, local.CenterDB, api.logger)
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
