package team_pk

import (
	"data_backend/apps/v2/internal/activity/team_pk/api"
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
		rg := r.Group("team-pk")
		teamPKApi := api.NewTeamPKApi()
		rg.GET("", local.PermissionGate.CheckPerm("activity_team_pk_view"), teamPKApi.List)
		rg.POST("/export", local.PermissionGate.CheckPerm("activity_team_pk_view"), teamPKApi.Export)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
