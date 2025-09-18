package invite

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/invite/api"
	"data_backend/apps/v2/internal/report/invite/dao"
	"data_backend/apps/v2/internal/report/invite/job"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{
			"0 5 * * ?": {job.NewInviteBetJob(), job.NewInviteBetDailyJob()},
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
		// invite-bet路由组
		rg := r.Group("invite-bet")
		// 初始化InviteBetApi实例，用于处理邀请的API请求
		inviteBetApi := api.NewInviteBetApi()
		rg.POST("generate", local.PermissionGate.CheckPerm("report_invite_generate"), inviteBetApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("report_invite_view"), inviteBetApi.List)
		rg.POST("export", local.PermissionGate.CheckPerm("report_invite_view"), inviteBetApi.Export)
	}

	{
		rg := r.Group("invite-bet/daily")
		InviteBetDailyApi := api.NewInviteBetDailyApi()
		rg.POST("generate", local.PermissionGate.CheckPerm("report_invite_daily_generate"), InviteBetDailyApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("report_invite_daily_view"), InviteBetDailyApi.List)
		rg.POST("export", local.PermissionGate.CheckPerm("report_invite_daily_view"), InviteBetDailyApi.Export)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{
		dao.InviteBet{},
		dao.InviteBetDaily{},
	}...)
}
