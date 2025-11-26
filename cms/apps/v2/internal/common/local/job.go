package local

import (
	"context"
	"time"

	"data_backend/pkg/cronjob"
	"data_backend/pkg/logger"
)

const JobRKey = "v2Jobs"

var JobWorker *cronjob.JobWorker
var JobWorkerLogger *logger.Logger
var (
	QueueCronChain cronjob.CronChain
	CronChain      cronjob.CronChain
)

func SetupJobs() (err error) {
	ctx := context.WithValue(Ctx, logger.ModuleKey, Module.Add(".job"))
	JobWorkerLogger = Logger.WithContext(ctx)
	alarm := NewAlarm(JobWorkerLogger)
	cronMessage := cronjob.NewJobConfig(ctx, JobWorkerLogger, alarm)

	CronChain = cronjob.NewCronChain([]cronjob.Wrapper{
		cronjob.RecoverWrapper(cronMessage),
		cronjob.SkipIfStillRunningWrapper(cronMessage),
		cronjob.LoggerWrapper(cronMessage),
		cronjob.TimeoutReminderWrapper(cronMessage, time.Minute),
	}...)
	QueueCronChain = cronjob.NewCronChain([]cronjob.Wrapper{
		cronjob.RecoverWrapper(cronMessage),
		cronjob.SkipIfStillRunningWrapper(cronMessage),
		cronjob.LoggerWrapper(cronMessage),
	}...)

	JobWorker = cronjob.NewJobWorker(ctx, JobRKey, cronMessage, RedisClient.Client)

	return nil
}
