package cost_award

import (
	"data_backend/apps/v2/internal/activity/cost_award/api"
	"data_backend/apps/v2/internal/activity/cost_award/dao"
	"data_backend/apps/v2/internal/activity/cost_award/job"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{
			"0 4 * * ?": {job.NewCostAwardJob()},
		},
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
		// cost-award路由组
		rg := r.Group("cost-award")
		costAwardApi := api.NewCostAwardApi()
		rg.GET("/generate", local.PermissionGate.CheckPerm("activity_cost_award_generate"), costAwardApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("activity_cost_award_view"), costAwardApi.List)
		rg.POST("/export", local.PermissionGate.CheckPerm("activity_cost_award_view"), costAwardApi.Export)
	}
	{
		// cost-award-log路由组
		rg := r.Group("cost-award-log")
		costAwardLogApi := api.NewCostAwardLogApi()
		rg.GET("/log-type/options", local.PermissionGate.CheckPerm("activity_cost_award_log_view"), costAwardLogApi.OptionsLogType)
		rg.GET("", local.PermissionGate.CheckPerm("activity_cost_award_log_view"), costAwardLogApi.List)
		rg.POST("/export", local.PermissionGate.CheckPerm("activity_cost_award_log_view"), costAwardLogApi.Export)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{
		dao.CostAward{},
	}...)
}
