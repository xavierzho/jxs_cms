package cohort

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/cohort/api"
	"data_backend/apps/v2/internal/report/cohort/dao"
	"data_backend/apps/v2/internal/report/cohort/job"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{
			"0 */2 * * ?": {job.NewCohortJob()},
			"0 3 * * ?":   {job.NewYesterdayCohortJob()},
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
		rg := r.Group("cohort")
		cohortApi := api.NewCohortApi()
		rg.POST("generate", local.PermissionGate.CheckPerm("report_cohort_generate"), cohortApi.Generate)
		rg.GET("", local.PermissionGate.CheckPerm("report_cohort_view"), cohortApi.All)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{
		dao.Cohort{},
	}...)
}
