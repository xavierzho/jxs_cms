package bet

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/bet/api"
	"data_backend/apps/v2/internal/report/bet/dao"
	"data_backend/apps/v2/internal/report/bet/job"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{
			"0 4 * * ?": {job.NewBetJob()},
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
		rg := r.Group("bet")
		betApi := api.NewBetApi()
		rg.POST("generate", local.PermissionGate.CheckPerm("report_bet_generate"), betApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("report_bet_view"), betApi.All)

	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{
		dao.Bet{},
	}...)
}
