package redemption_code

import (
	"data_backend/apps/v2/internal/activity/redemption_code/api"
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

func InitRouter(r *gin.RouterGroup) (err error) {
	{
		rg := r.Group("redemption-code")
		redemptionCodeApi := api.NewRedemptionCodeApi()
		rg.GET("", local.PermissionGate.CheckPerm("activity_redemption_code_view"), redemptionCodeApi.Log)
		rg.GET("/award-detail", local.PermissionGate.CheckPerm("activity_redemption_code_view"), redemptionCodeApi.GetAwardDetail)
		rg.POST("/export", local.PermissionGate.CheckPerm("activity_redemption_code_view"), redemptionCodeApi.ExportLog)
		rg.POST("/detail/export", local.PermissionGate.CheckPerm("activity_redemption_code_view"), redemptionCodeApi.ExportAwardDetail)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
