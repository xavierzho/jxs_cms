package job

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/order/form"
	"data_backend/apps/v2/internal/report/order/service"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

// InviteBetJob 结构体用于邀请投注任务。
type DeliveryOrderJob struct {
	ctx    *gin.Context
	logger *logger.Logger
	alarm  message.Alarm
}

// NewInviteBetJob 创建并初始化一个新的 InviteBetJob 实例
func NewDeliveryOrderJob() *DeliveryOrderJob {
	// 设置日志模块和上下文
	log := local.JobWorkerLogger.WithContext(context.WithValue(local.JobWorkerLogger.Context, logger.ModuleKey, local.JobWorkerLogger.ModuleKey().Add(".DeliveryOrderJob")))
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Request = ctx.Request.WithContext(log.Context)
	return &DeliveryOrderJob{
		ctx:    ctx,
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

// Name 返回作业的名称
func (*DeliveryOrderJob) Name() string {
	return "DeliveryOrderJob"
}

// Run 将当前作业实例添加到作业队列中
func (j *DeliveryOrderJob) Run() {
	local.JobWorker.AddJobToQueue(j.Name())
}

// Work 执行作业的具体逻辑
func (j *DeliveryOrderJob) Work() {
	// 获取前一天的日期字符串
	cDateStr := time.Now().AddDate(0, 0, -1).Format(pkg.DATE_FORMAT)
	// 初始化成本奖励服务
	svc := service.NewDeliveryOrderSvc(j.ctx, local.CMSDB, local.CenterDB, j.logger)
	if err := svc.Generate(&form.GenerateRequest{DateRangeRequest: iForm.DateRangeRequest{DateRange: [2]string{cDateStr, cDateStr}}}); err != nil {
		j.alarm.AlertErrorMsg(fmt.Sprintf("DeliveryOrderSvc.Generate error: %v", err), message.CMS_ID)
	}
}
