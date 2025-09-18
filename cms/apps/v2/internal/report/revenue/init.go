package revenue

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/revenue/api"
	"data_backend/apps/v2/internal/report/revenue/dao"
	"data_backend/apps/v2/internal/report/revenue/job"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{
			"0 0 * * ?": {job.NewRevenueBalanceJob()},
			"0 3 * * ?": {job.NewYesterdayRevenueJob()},
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
		rg := r.Group("revenue")
		revenueApi := api.NewRevenueApi()
		rg.POST("generate", local.PermissionGate.CheckPerm("report_revenue_generate"), revenueApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("report_revenue_view"), revenueApi.All)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{
		dao.Active{}, dao.Pating{}, dao.Wastage{},
		dao.Balance{}, dao.Draw{}, dao.Pay{},
	}...)
}
