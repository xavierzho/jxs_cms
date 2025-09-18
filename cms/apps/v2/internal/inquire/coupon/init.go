package coupon

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/coupon/api"
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
		rg := r.Group("coupon")
		couponApi := api.NewCouponApi()
		rg.GET("/type/options", local.PermissionGate.CheckPerm("inquire_coupon_view"), couponApi.OptionsCouponType)
		rg.GET("/action/options", local.PermissionGate.CheckPerm("inquire_coupon_view"), couponApi.OptionsCouponActionType)
		rg.GET("", local.PermissionGate.CheckPerm("inquire_coupon_view"), couponApi.List)
		rg.POST("/export", local.PermissionGate.CheckPerm("inquire_coupon_view"), couponApi.Export)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
