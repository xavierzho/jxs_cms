package realtime

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/realtime/api"
	"data_backend/apps/v2/internal/report/realtime/job"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{
			"*/10 * * * ?": {job.NewRealtimeJob()},
		},
	)
}

func AddQueueJob() error {
	return local.QueueWorker.AddQueueJob(
		[]*queue.QueueJob{},
	)
}

func InitRouter(r *gin.RouterGroup) (err error) {
	{
		rg := r.Group("realtime")
		realtimeApi := api.NewRealtimeApi()
		rg.POST("cached", local.PermissionGate.CheckPerm("report_realtime_cached"), realtimeApi.Cached)
		rg.GET("", local.PermissionGate.CheckPerm("report_realtime_view"), realtimeApi.All)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
