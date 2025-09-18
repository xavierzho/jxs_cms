package job

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/market/form"
	"data_backend/apps/v2/internal/report/market/service"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type MarketJob struct {
	ctx    *gin.Context
	logger *logger.Logger
	alarm  message.Alarm
}

func NewMarketJob() *MarketJob {
	log := local.JobWorkerLogger.WithContext(context.WithValue(local.JobWorkerLogger.Context, logger.ModuleKey, local.JobWorkerLogger.ModuleKey().Add(".MarketJob")))
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Request = ctx.Request.WithContext(log.Context)
	return &MarketJob{
		ctx:    ctx,
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (*MarketJob) Name() string {
	return "MarketJob"
}

func (j *MarketJob) Run() {
	local.JobWorker.AddJobToQueue(j.Name())
}

func (j *MarketJob) Work() {
	cDateStr := time.Now().AddDate(0, 0, -1).Format(pkg.DATE_FORMAT)
	svc := service.NewMarketSvc(j.ctx, local.CMSDB, local.CenterDB, j.logger)
	if err := svc.Generate(&form.GenerateRequest{DateRangeRequest: iForm.DateRangeRequest{DateRange: [2]string{cDateStr, cDateStr}}}); err != nil {
		j.alarm.AlertErrorMsg(fmt.Sprintf("MarketSvc.Generate: %v", err), message.CMS_ID)
	}
}
