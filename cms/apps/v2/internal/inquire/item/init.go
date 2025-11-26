package item

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/item/api"
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
		rg := r.Group("item")
		itemApi := api.NewItemApi()
		rg.GET("/log-type/options", local.PermissionGate.CheckPerm("inquire_item_view"), itemApi.OptionsLogType)
		rg.GET("/log", local.PermissionGate.CheckPerm("inquire_item_view"), itemApi.GetLog)
		rg.POST("/log/export", local.PermissionGate.CheckPerm("inquire_item_view"), itemApi.ExportLog)
		rg.GET("/detail", local.PermissionGate.CheckPerm("inquire_item_view"), itemApi.GetDetail)
		rg.POST("/detail/export", local.PermissionGate.CheckPerm("inquire_item_view"), itemApi.ExportDetail)
		rg.GET("/bet", local.PermissionGate.CheckPerm("inquire_bet_item_view"), itemApi.ListBetDetail)
		rg.POST("/bet/export", local.PermissionGate.CheckPerm("inquire_bet_item_view"), itemApi.ExportBetDetail)
		rg.GET("/revenue", local.PermissionGate.CheckPerm("inquire_revenue_item_view"), itemApi.GetRevenue)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
