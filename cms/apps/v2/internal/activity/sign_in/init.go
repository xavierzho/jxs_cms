package sign_in

import (
	"data_backend/apps/v2/internal/activity/sign_in/api"
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
		rg := r.Group("sign-in")
		SignInApi := api.NewSignInApi()
		rg.GET("", local.PermissionGate.CheckPerm("activity_sign_in_view"), SignInApi.List)
		rg.POST("/export", local.PermissionGate.CheckPerm("activity_sign_in_view"), SignInApi.Export)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
