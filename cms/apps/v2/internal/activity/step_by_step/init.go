package step_by_step

import (
	"data_backend/apps/v2/internal/activity/step_by_step/api"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{},
	)
}

func AddQueueJob() error {
	return local.QueueWorker.AddQueueJob(
		[]*queue.QueueJob{},
	)
}

// 活动路由
func InitRouter(r *gin.RouterGroup) (err error) {
	{
		rg := r.Group("step-by-step")
		stepByStepApi := api.NewStepByStepApi()
		rg.GET("", local.PermissionGate.CheckPerm("activity_step_by_step_view"), stepByStepApi.LogList)
		rg.POST("/export", local.PermissionGate.CheckPerm("activity_step_by_step_view"), stepByStepApi.LogExport)
		rg.GET("/detail", local.PermissionGate.CheckPerm("activity_step_by_step_view"), stepByStepApi.Detail)
		rg.POST("/detail/export", local.PermissionGate.CheckPerm("activity_step_by_step_view"), stepByStepApi.DetailExport)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
