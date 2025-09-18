package market

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/market/api"
	"data_backend/apps/v2/internal/report/market/dao"
	"data_backend/apps/v2/internal/report/market/job"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{
			"0 4 * * ?": {job.NewMarketJob()},
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
		rg := r.Group("market")
		marketApi := api.NewMarketApi()
		rg.POST("generate", local.PermissionGate.CheckPerm("report_market_generate"), marketApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("report_market_view"), marketApi.List)
		rg.POST("export", local.PermissionGate.CheckPerm("report_market_view"), marketApi.Export)

	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{
		dao.Market{},
	}...)
}
