package job

import (
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/pkg"
)

type YesterdayCohortJob struct {
	*CohortJob
}

func NewYesterdayCohortJob() *YesterdayCohortJob {
	return &YesterdayCohortJob{
		NewCohortJob(),
	}
}

func (*YesterdayCohortJob) Name() string {
	return "YesterdayCohortJob"
}

func (j *YesterdayCohortJob) Run() {
	local.JobWorker.AddJobToQueue(j.Name())
}

func (j *YesterdayCohortJob) Work() {
	// 更新昨天的cohort
	j.now = time.Now().AddDate(0, 0, -1).Format(pkg.DATE_FORMAT)
	j.CohortJob.Work()
}
