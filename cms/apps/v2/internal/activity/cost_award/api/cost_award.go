package api

import (
	"context"
	"fmt"

	"data_backend/apps/v2/internal/activity/cost_award/form"
	"data_backend/apps/v2/internal/activity/cost_award/service"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/app"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type CostAwardApi struct {
	logger *logger.Logger //记录api操作日志
	alarm  message.Alarm  //记录报警信息
}

func NewCostAwardApi() *CostAwardApi {
	// 生成与当前模块相关的日志记录器实例，用于后续的日志记录。
	log := local.Logger.WithContext(context.WithValue(local.Ctx, logger.ModuleKey, local.Module.Add(".CostAwardApi")))
	// 初始化并返回 CostAwardApi 实例，为其提供日志记录器和新的报警器实例。
	return &CostAwardApi{
		logger: log,
		alarm:  local.NewAlarm(log),
	}

}

func (api *CostAwardApi) Generate(ctx *gin.Context) {
	params := &form.GenerateRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewCostAwardSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	go func() {
		if err := svc.Generate(params); err != nil {
			api.alarm.AlertErrorMsg(fmt.Sprintf("CostAwardSvc.Generate: %v", err), message.CMS_ID)
		}
	}()

	response.ToResponseOK()
}

func (api *CostAwardApi) List(ctx *gin.Context) {
	params := &form.ListRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}

	svc := service.NewCostAwardSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
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

// 根据时间筛选导出到Excel
func (api *CostAwardApi) Export(ctx *gin.Context) {
	//初始化请求参数
	params := &form.AllRequest{}
	response := app.NewResponse(ctx)
	if ok := response.BindAndValid(ctx, params, api.logger); !ok {
		return
	}
	//创建服务对象
	svc := service.NewCostAwardSvc(ctx, local.CMSDB, local.CenterDB, api.logger)
	excelModel, err := svc.Export(params) //server层执行导出操作
	if err != nil {
		api.logger.Errorf("Export: %v", err)
		response.ToErrorResponse(err)
		return
	}
	//导出文件
	e := response.ExportFile(ctx, excelModel.Excelize, excelModel.FileName)
	if e != nil {
		api.logger.Errorf("response.ExportFile err: %v", e.Error())
	}
}
