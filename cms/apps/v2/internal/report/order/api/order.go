package api

import (
	"context"
	"fmt"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/order/form"
	"data_backend/apps/v2/internal/report/order/service"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type DeliveryOrderApi struct {
	logger *logger.Logger //记录日志
	alarm  message.Alarm  //预警信息
}

// NewInviteBetApi 创建并返回一个新的 InviteBetApi 实例
func NewDeliveryOrderApi() *DeliveryOrderApi {
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".DeliveryOrderApi")))
	return &DeliveryOrderApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}

}

func (api *DeliveryOrderApi) Generate(ctx *gin.Context) {
	params := &form.GenerateRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewDeliveryOrderSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	go func() {
		if err := svc.Generate(params); err != nil {
			api.alarm.AlertErrorMsg(fmt.Sprintf("DeliveryOrderSvc.Generate: %v", err), message.CMS_ID)
		}
	}()

	response.ToResponseOK()
}

func (api *DeliveryOrderApi) List(ctx *gin.Context) {
	params := &form.ListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewDeliveryOrderSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
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

func (api *DeliveryOrderApi) Export(ctx *gin.Context) {
	params := &form.AllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}
	svc := service.NewDeliveryOrderSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
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
