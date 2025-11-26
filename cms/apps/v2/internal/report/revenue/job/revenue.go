package job

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/revenue/form"
	"data_backend/apps/v2/internal/report/revenue/service"
	"data_backend/pkg"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

// 更新当天的数据
type RevenueBalanceJob struct {
	ctx    *gin.Context
	logger *logger.Logger
	alarm  message.Alarm
}

func NewRevenueBalanceJob() *RevenueBalanceJob {
	log := local.JobWorkerLogger.WithContext(context.WithValue(local.JobWorkerLogger.Context, logger.ModuleKey, local.JobWorkerLogger.ModuleKey().Add(".RevenueBalanceJob")))
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Request = ctx.Request.WithContext(log.Context)
	return &RevenueBalanceJob{
		ctx:    ctx,
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (*RevenueBalanceJob) Name() string {
	return "RevenueBalanceJob"
}

func (j *RevenueBalanceJob) Run() {
	local.JobWorker.AddJobToQueue(j.Name())
}

func (j *RevenueBalanceJob) Work() {
	cDateStr := time.Now().AddDate(0, 0, -1).Format(pkg.DATE_FORMAT)
	svc := service.NewRevenueSvc(j.ctx, local.CMSDB, local.CenterDB, j.logger)

	// 更新昨天的 支付, 退款(￥), 钱包, 活跃, 参与数据
	err := svc.Generate(&form.GenerateRequest{
		DateRange: [2]string{cDateStr, cDateStr},
		DataTypeList: []string{
			form.REVENUE_DATA_TYPE_BALANCE,
		},
	})
	if err != nil {
		j.alarm.AlertErrorMsg(fmt.Sprintf("RevenueSvc.Generate: %v", err), message.CmsId)
	}
}
