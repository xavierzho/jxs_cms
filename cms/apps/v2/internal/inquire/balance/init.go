package balance

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/balance/api"
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
		rg := r.Group("balance")
		balanceApi := api.NewBalanceApi()
		rg.GET("/source-type/options", local.PermissionGate.CheckPerm("inquire_balance_view"), balanceApi.OptionsSourceType)
		rg.GET("/channel-type/options", local.PermissionGate.CheckPerm("inquire_balance_view"), balanceApi.OptionsChannelType)
		rg.GET("/pay-source-type/options", local.PermissionGate.CheckPerm("inquire_balance_view"), balanceApi.OptionsPaySourceType)
		rg.GET("/balance-type/options", local.PermissionGate.CheckPerm("inquire_balance_view"), balanceApi.OptionsBalanceType)

		rg.GET("", local.PermissionGate.CheckPerm("inquire_balance_view"), balanceApi.List)
		rg.POST("comment", local.PermissionGate.CheckPerm("inquire_balance_view"), balanceApi.AddComment)
		rg.DELETE("comment", local.PermissionGate.CheckPerm("inquire_balance_view"), balanceApi.DeleteComment)
		rg.POST("/export", local.PermissionGate.CheckPerm("inquire_balance_view"), balanceApi.Export)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
