package job

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"data_backend/apps/v2/internal/activity/cost_award/form"
	"data_backend/apps/v2/internal/activity/cost_award/service"
	"data_backend/apps/v2/internal/common/local"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type CostAwardJob struct {
	ctx    *gin.Context
	logger *logger.Logger
	alarm  message.Alarm
}

func NewCostAwardJob() *CostAwardJob {
	log := local.JobWorkerLogger.WithContext(context.WithValue(local.JobWorkerLogger.Context, logger.ModuleKey, local.JobWorkerLogger.ModuleKey().Add(".CostAwardJob")))
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Request = ctx.Request.WithContext(log.Context)
	return &CostAwardJob{
		ctx:    ctx,
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (*CostAwardJob) Name() string {
	return "CostAwardJob"
}

func (j *CostAwardJob) Run() {
	local.JobWorker.AddJobToQueue(j.Name())
}

func (j *CostAwardJob) Work() {
	cDateStr := time.Now().AddDate(0, 0, -1).Format(pkg.DATE_FORMAT)
	svc := service.NewCostAwardSvc(j.ctx, local.CMSDB, local.CenterDB, j.logger)
	if err := svc.Generate(&form.GenerateRequest{DateRangeRequest: iForm.DateRangeRequest{DateRange: [2]string{cDateStr, cDateStr}}}); err != nil {
		j.alarm.AlertErrorMsg(fmt.Sprintf("CostAwardSvc.Generate: %v", err), message.CMS_ID)
	}
}
