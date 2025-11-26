package recall

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/recall/api"
	"data_backend/apps/v2/internal/report/recall/dao"
	"data_backend/apps/v2/internal/report/recall/job"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{
			"0 5 * * ?": {job.NewRecallJob(), job.NewRecallDailyJob()},
		},
	)
}

func AddQueueJob() error {
	return local.QueueWorker.AddQueueJob(
		[]*queue.QueueJob{},
	)
}

// 邀请路由
func InitRouter(r *gin.RouterGroup) (err error) {
	{
		// recall路由组
		rg := r.Group("recall")
		// 初始化NewRecallApi实例，用于处理邀请的API请求
		recallApi := api.NewRecallApi()
		rg.POST("generate", local.PermissionGate.CheckPerm("report_recall_generate"), recallApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("report_recall_view"), recallApi.List)
		rg.POST("export", local.PermissionGate.CheckPerm("report_recall_view"), recallApi.Export)
	}

	{
		rg := r.Group("recall/daily")
		recallDailyApi := api.NewRecallDailyApi()
		rg.POST("generate", local.PermissionGate.CheckPerm("report_recall_daily_generate"), recallDailyApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("report_recall_daily_view"), recallDailyApi.List)
		rg.POST("export", local.PermissionGate.CheckPerm("report_recall_daily_view"), recallDailyApi.Export)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{
		dao.Recall{},
		dao.RecallDaily{},
	}...)
}
