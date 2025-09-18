package gacha

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/gacha/api"
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
		rg := r.Group("gacha")
		gachaApi := api.NewGachaApi()
		rg.GET("/type/options", local.PermissionGate.CheckPerm("inquire_gacha_view"), gachaApi.OptionsGachaType)
		rg.GET("/revenue", local.PermissionGate.CheckPerm("inquire_gacha_view"), gachaApi.GetRevenue)
		rg.GET("/detail", local.PermissionGate.CheckPerm("inquire_gacha_view"), gachaApi.GetDetail)
		rg.POST("/detail/export", local.PermissionGate.CheckPerm("inquire_gacha_view"), gachaApi.ExportDetail)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
