package marketing

import (
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/marketing/api"
	"data_backend/apps/v2/internal/marketing/dao"
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

// InitRouter registers the marketing routes
func InitRouter(r *gin.RouterGroup) (err error) {
	{
		rg := r.Group("marketing")
		// Note: We might want to add JWT middleware here if this is only for logged-in users
		// rg.Use(local.JWT.JWT())

		marketingAPI := api.NewMarketingAPI(local.CMSDB, local.Logger)

		// Endpoint to record attribution (e.g. called by App after login/install)
		rg.POST("attribution", marketingAPI.RecordAttribution)
	}

	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{
		dao.UserAttribution{},
	}...)
}
