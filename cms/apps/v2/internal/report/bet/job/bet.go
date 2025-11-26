package job

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/bet/form"
	"data_backend/apps/v2/internal/report/bet/service"
	"data_backend/pkg"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type BetJob struct {
	ctx    *gin.Context
	logger *logger.Logger
	alarm  message.Alarm
}

func NewBetJob() *BetJob {
	log := local.JobWorkerLogger.WithContext(context.WithValue(local.JobWorkerLogger.Context, logger.ModuleKey, local.JobWorkerLogger.ModuleKey().Add(".BetJob")))
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Request = ctx.Request.WithContext(log.Context)
	return &BetJob{
		ctx:    ctx,
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (*BetJob) Name() string {
	return "BetJob"
}

func (j *BetJob) Run() {
	local.JobWorker.AddJobToQueue(j.Name())
}

func (j *BetJob) Work() {
	cDateStr := time.Now().AddDate(0, 0, -1).Format(pkg.DATE_FORMAT)
	svc := service.NewBetSvc(j.ctx, local.CMSDB, local.CenterDB, j.logger)
	if err := svc.Generate(&form.GenerateRequest{DateRange: [2]string{cDateStr, cDateStr}}); err != nil {
		j.alarm.AlertErrorMsg(fmt.Sprintf("BetSvc.Generate: %v", err), message.CmsId)
	}
}
