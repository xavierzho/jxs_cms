package order

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/order/api"
	"data_backend/apps/v2/internal/report/order/dao"
	"data_backend/apps/v2/internal/report/order/job"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{
			"0 5 * * ?": {job.NewDeliveryOrderJob()},
		},
	)
}

func AddQueueJob() error {
	return local.QueueWorker.AddQueueJob(
		[]*queue.QueueJob{},
	)
}

// 发货报表路由
func InitRouter(r *gin.RouterGroup) (err error) {
	{
		rg := r.Group("delivery_order")
		DeliveryOrderDailyApi := api.NewDeliveryOrderApi()
		rg.POST("generate", local.PermissionGate.CheckPerm("report_order_daily_generate"), DeliveryOrderDailyApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("report_order_view"), DeliveryOrderDailyApi.List)
		rg.POST("export", local.PermissionGate.CheckPerm("report_order_view"), DeliveryOrderDailyApi.Export)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{
		dao.DeliveryOrder{},
	}...)
}
