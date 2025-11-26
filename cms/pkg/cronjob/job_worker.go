package cronjob

import (
	"context"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// JobWorker manages the lifecycle of cron jobs and queue-based jobs.
// It maintains job registries, handles scheduling, and coordinates between
// scheduled execution and queue-based execution.
type JobWorker struct {
	ctx               context.Context
	queueRKey         string
	config            *JobConfig
	redisClient       *redis.Client
	jobList           map[string][]CronCommonJob
	frequentlyJobList map[string][]CronCommonJob
	queueJobDict      map[string]CronQueueJob
	mu                sync.RWMutex // Protects queueJobDict
	crons             []*cron.Cron // Track all cron instances for cleanup
}

// NewJobWorker creates a new JobWorker instance.
//
// Parameters:
//   - ctx: Context for job execution
//   - queueRKey: Redis key for the job queue
//   - config: Job configuration containing logger and alarm
//   - redisClient: Redis client for queue operations
func NewJobWorker(ctx context.Context, queueRKey string, config *JobConfig, redisClient *redis.Client) *JobWorker {
	return &JobWorker{
		ctx:               ctx,
		queueRKey:         queueRKey,
		config:            config,
		redisClient:       redisClient,
		jobList:           make(map[string][]CronCommonJob),
		frequentlyJobList: make(map[string][]CronCommonJob),
		queueJobDict:      make(map[string]CronQueueJob),
		crons:             make([]*cron.Cron, 0),
		mu:                sync.RWMutex{},
	}
}

// QueueRKey returns the Redis key used for the job queue.
func (j *JobWorker) QueueRKey() string {
	return j.queueRKey
}

// addQueueJobDict adds a queue job to the internal registry.
// Returns an error if a job with the same name already exists.
func (j *JobWorker) addQueueJobDict(job CronQueueJob) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	if _, ok := j.queueJobDict[job.Name()]; ok {
		return NewErrJobAlreadyExists(job.Name())
	}
	j.queueJobDict[job.Name()] = job
	return nil
}

// AddJobList adds jobs to the regular job list with their cron schedules.
//
// Parameters:
//   - jobList: Map of cron schedule expressions to lists of jobs
//
// Example:
//
//	worker.AddJobList(map[string][]CronCommonJob{
//	  "0 */6 * * *": {dailyReportJob},  // Every 6 hours
//	  "0 0 * * *":   {cleanupJob},      // Daily at midnight
//	})
func (j *JobWorker) AddJobList(jobList map[string][]CronCommonJob) error {
	for spec, jobs := range jobList {
		j.jobList[spec] = append(j.jobList[spec], jobs...)
	}
	return nil
}

// StartJob starts the regular cron scheduler with the provided wrapper chains.
//
// Parameters:
//   - cronChain: Wrapper chain for regular jobs
//   - queueCronChain: Wrapper chain for queue-based jobs
//
// Returns nil if no jobs are registered.
func (j *JobWorker) StartJob(cronChain, queueCronChain CronChain) error {
	if len(j.jobList) == 0 {
		return nil
	}
	c := cron.New()
	j.mu.Lock()
	j.crons = append(j.crons, c)
	j.mu.Unlock()
	return j.addJob(j.jobList, c, cronChain, queueCronChain)
}

// AddFrequentlyJobList adds jobs that need to run with second-level precision.
//
// Parameters:
//   - jobList: Map of cron schedule expressions (with seconds) to lists of jobs
//
// Example:
//
//	worker.AddFrequentlyJobList(map[string][]CronCommonJob{
//	  "*/30 * * * * *": {healthCheckJob},  // Every 30 seconds
//	})
func (j *JobWorker) AddFrequentlyJobList(jobList map[string][]CronCommonJob) error {
	for spec, jobs := range jobList {
		j.frequentlyJobList[spec] = append(j.frequentlyJobList[spec], jobs...)
	}
	return nil
}

// StartFrequentlyJob starts the cron scheduler for jobs requiring second-level precision.
//
// Parameters:
//   - cronChain: Wrapper chain to apply to all frequently-run jobs
func (j *JobWorker) StartFrequentlyJob(cronChain CronChain) error {
	if len(j.frequentlyJobList) == 0 {
		return nil
	}
	c := cron.New(cron.WithSeconds())
	j.mu.Lock()
	j.crons = append(j.crons, c)
	j.mu.Unlock()
	return j.addJob(j.frequentlyJobList, c, cronChain, cronChain)
}

// addJob is an internal method that adds jobs to a cron instance with appropriate wrappers.
func (j *JobWorker) addJob(jobList map[string][]CronCommonJob, c *cron.Cron, cronChain CronChain, queueCronChain CronChain) error {
	for spec, jobs := range jobList {
		for _, job := range jobs {
			var wrappedJob CronCommonJob

			// Apply appropriate wrapper chain based on job type
			if qJob, ok := job.(CronQueueJob); ok {
				wrappedJob = queueCronChain.then(job)
				if err := j.addQueueJobDict(qJob); err != nil {
					c.Stop()
					return err
				}
			} else {
				wrappedJob = cronChain.then(job)
			}

			// Schedule the job
			if _, err := c.AddJob(spec, wrappedJob); err != nil {
				c.Stop()
				return NewErrJobScheduleFailed(job.Name(), err)
			}
		}
	}

	c.Start()
	return nil
}

// AddJobToQueue pushes a job name to the Redis queue for asynchronous execution.
// This triggers the execution of the corresponding CronQueueJob's Work() method.
//
// Parameters:
//   - jobName: The name of the job to execute
//
// Returns an error if the job name is empty or if the Redis operation fails.
func (j *JobWorker) AddJobToQueue(jobName string) error {
	if strings.TrimSpace(jobName) == "" {
		err := NewErrEmptyJobName()
		j.config.Alarm(logrus.ErrorLevel, err.Error())
		return err
	}

	if _, err := j.redisClient.RPush(j.ctx, j.queueRKey, jobName).Result(); err != nil {
		wrappedErr := NewErrQueuePushFailed(jobName, err)
		j.config.Alarm(logrus.ErrorLevel, wrappedErr.Error())
		return wrappedErr
	}

	return nil
}

// QueueJobRun executes the Work() method of a queue job by name.
// This is typically called by a queue consumer that pops job names from Redis.
//
// Parameters:
//   - jobName: The name of the job to execute
//
// Returns an error if the job is not found in the queue job registry.
func (j *JobWorker) QueueJobRun(jobName string) error {
	j.mu.RLock()
	job, ok := j.queueJobDict[jobName]
	j.mu.RUnlock()

	if !ok {
		err := NewErrJobNotFound(jobName)
		j.config.Logf(logrus.ErrorLevel, err.Error())
		return err
	}

	job.Work()
	return nil
}

// Stop gracefully stops all running cron schedulers.
// This should be called during application shutdown to ensure clean termination.
func (j *JobWorker) Stop() {
	j.mu.RLock()
	defer j.mu.RUnlock()

	for _, c := range j.crons {
		if c != nil {
			c.Stop()
		}
	}
}
