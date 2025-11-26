package item

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/task/api"
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

func InitRouter(r *gin.RouterGroup) (err error) {
	{
		rg := r.Group("task")
		taskApi := api.NewTaskApi()
		rg.GET("/type/options", local.PermissionGate.CheckPerm("inquire_task_view"), taskApi.OptionsType)
		rg.GET("/key/options", local.PermissionGate.CheckPerm("inquire_task_view"), taskApi.OptionsKey)
		rg.GET("", local.PermissionGate.CheckPerm("inquire_task_view"), taskApi.List)
		rg.GET("/award-detail", local.PermissionGate.CheckPerm("inquire_task_view"), taskApi.GetAwardDetail)
		rg.POST("/export", local.PermissionGate.CheckPerm("inquire_task_view"), taskApi.ExportList)
		rg.POST("/detail/export", local.PermissionGate.CheckPerm("inquire_task_view"), taskApi.ExportAwardDetail)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
