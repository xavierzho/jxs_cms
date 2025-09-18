package job

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/dashboard/form"
	"data_backend/apps/v2/internal/report/dashboard/service"
	"data_backend/pkg"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type DashboardJob struct {
	ctx    *gin.Context
	logger *logger.Logger
	alarm  message.Alarm
}

func NewDashboardJob() *DashboardJob {
	log := local.JobWorkerLogger.WithContext(context.WithValue(local.JobWorkerLogger.Context, logger.ModuleKey, local.JobWorkerLogger.ModuleKey().Add(".DashboardJob")))
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Request = ctx.Request.WithContext(log.Context)
	return &DashboardJob{
		ctx:    ctx,
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (*DashboardJob) Name() string {
	return "DashboardJob"
}

func (j *DashboardJob) Run() {
	local.JobWorker.AddJobToQueue(j.Name())
}

func (j *DashboardJob) Work() {
	cDateStr := time.Now().AddDate(0, 0, -1).Format(pkg.DATE_FORMAT)
	svc := service.NewDashboardSvc(j.ctx, local.CMSDB, local.CenterDB, j.logger)
	if err := svc.Generate(&form.GenerateRequest{DateRange: [2]string{cDateStr, cDateStr}}); err != nil {
		j.alarm.AlertErrorMsg(fmt.Sprintf("DashboardSvc.Generate: %v", err), message.CMS_ID)
	}
}
