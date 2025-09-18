package invite

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/invite/api"
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

// 路由组
func InitRouter(r *gin.RouterGroup) (err error) {
	{
		rg := r.Group("invite-rec")
		InviteRecApi := api.NewInviteRecApi() // 初始化InviteRecApi实例，用于处理邀请的API请求
		rg.GET("", local.PermissionGate.CheckPerm("inquire_invite_view"), InviteRecApi.List)
		rg.POST("/export", local.PermissionGate.CheckPerm("inquire_invite_view"), InviteRecApi.Export)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
