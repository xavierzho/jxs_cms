/*
Package cronjob provides an enhanced wrapper around github.com/robfig/cron/v3
with additional features for running and managing scheduled jobs.

Key features:
 1. Extended job interfaces (CronCommonJob and CronQueueJob)
 2. Custom logging integration with structured logging support
 3. Decorator/wrapper functionality (recovery, logging, timeout, skip-if-running)
 4. Redis-based job queue management

Example usage:

	ctx := context.Background()
	logger := myLogger // implements cronjob.Logger
	alarm := myAlarm   // implements cronjob.Alarm

	config := cronjob.NewJobConfig(ctx, logger, alarm)
	worker := cronjob.NewJobWorker(ctx, "myQueue", config, redisClient)

	// Add jobs
	worker.AddJobList(map[string][]cronjob.CronCommonJob{
		"0 * * * *": {myJob}, // Every hour
	})

	// Start with wrappers
	chain := cronjob.NewCronChain(
		cronjob.RecoverWrapper(config),
		cronjob.LoggerWrapper(config),
	)
	worker.StartJob(chain, chain)
*/
package cronjob
