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

type YesterdayRevenueJob struct {
	ctx    *gin.Context
	logger *logger.Logger
	alarm  message.Alarm
}

func NewYesterdayRevenueJob() *YesterdayRevenueJob {
	log := local.JobWorkerLogger.WithContext(context.WithValue(local.JobWorkerLogger.Context, logger.ModuleKey, local.JobWorkerLogger.ModuleKey().Add(".YesterdayRevenueJob")))
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Request = ctx.Request.WithContext(log.Context)
	return &YesterdayRevenueJob{
		ctx:    ctx,
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (*YesterdayRevenueJob) Name() string {
	return "YesterdayRevenueJob"
}

func (j *YesterdayRevenueJob) Run() {
	local.JobWorker.AddJobToQueue(j.Name())
}

func (j *YesterdayRevenueJob) Work() {
	now := time.Now()
	svc := service.NewRevenueSvc(j.ctx, local.CMSDB, local.CenterDB, j.logger)

	// 更新昨天的 支付, 退款(￥), 钱包, 活跃, 参与数据
	err := svc.Generate(&form.GenerateRequest{
		DateRange: [2]string{now.AddDate(0, 0, -1).Format(pkg.DATE_FORMAT), now.AddDate(0, 0, -1).Format(pkg.DATE_FORMAT)},
		DataTypeList: []string{
			form.REVENUE_DATA_TYPE_PAY,
			form.REVENUE_DATA_TYPE_DRAW,
			// form.REVENUE_DATA_TYPE_BALANCE,
			form.REVENUE_DATA_TYPE_ACTIVE,
			form.REVENUE_DATA_TYPE_PATING,
		},
	})
	if err != nil {
		j.alarm.AlertErrorMsg(fmt.Sprintf("RevenueSvc.Generate: %v", err), message.CmsId)
	}

	// 更新昨天的 支付, 退款(￥), 钱包, 活跃, 参与数据
	err = svc.Generate(&form.GenerateRequest{
		DateRange: [2]string{now.AddDate(0, 0, -7).Format(pkg.DATE_FORMAT), now.AddDate(0, 0, -7).Format(pkg.DATE_FORMAT)},
		DataTypeList: []string{
			form.REVENUE_DATA_TYPE_WASTAGE,
		},
	})
	if err != nil {
		j.alarm.AlertErrorMsg(fmt.Sprintf("RevenueSvc.Generate: %v", err), message.CmsId)
	}
}
