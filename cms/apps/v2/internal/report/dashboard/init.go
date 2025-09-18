package dashboard

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/dashboard/api"
	"data_backend/apps/v2/internal/report/dashboard/dao"
	"data_backend/apps/v2/internal/report/dashboard/job"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{
			"0 4 * * ?": {job.NewDashboardJob()},
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
		rg := r.Group("dashboard")
		dashboardApi := api.NewDashboardApi()
		rg.POST("generate", local.PermissionGate.CheckPerm("report_dashboard_generate"), dashboardApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("report_dashboard_view"), dashboardApi.List)

	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{
		dao.Dashboard{},
	}...)
}
