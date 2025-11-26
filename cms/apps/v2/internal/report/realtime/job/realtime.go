package job

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/realtime/form"
	"data_backend/apps/v2/internal/report/realtime/service"
	"data_backend/pkg"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type RealtimeJob struct {
	ctx    *gin.Context
	logger *logger.Logger
	alarm  message.Alarm
}

func NewRealtimeJob() *RealtimeJob {
	log := local.JobWorkerLogger.WithContext(context.WithValue(local.JobWorkerLogger.Context, logger.ModuleKey, local.JobWorkerLogger.ModuleKey().Add(".RealtimeJob")))
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Request = ctx.Request.WithContext(log.Context)
	return &RealtimeJob{
		ctx:    ctx,
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (*RealtimeJob) Name() string {
	return "RealtimeJob"
}

func (j *RealtimeJob) Run() {
	local.JobWorker.AddJobToQueue(j.Name())
}

func (j *RealtimeJob) Work() {
	now := time.Now()

	svc := service.NewRealtimeSvc(j.ctx, local.CenterDB, local.RedisClient, j.logger)
	if err := svc.Cached(&form.CachedRequest{DateTime: now.Format(pkg.DATE_TIME_FORMAT)}); err != nil {
		j.alarm.AlertErrorMsg(fmt.Sprintf("RealtimeSvc.Cached: %v", err), message.CmsId)
	}
}
