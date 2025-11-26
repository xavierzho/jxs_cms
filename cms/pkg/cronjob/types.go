package cronjob

import (
	"context"

	"github.com/sirupsen/logrus"
)

// CronCommonJob represents a basic cron job that can be scheduled and executed.
// All cron jobs must implement this interface to provide their execution logic
// and a unique identifier.
type CronCommonJob interface {
	// Run executes the job's main logic. This method is called by the cron scheduler
	// when the job's schedule triggers.
	Run()

	// Name returns a unique identifier for this job. This is used for logging,
	// error reporting, and job queue management.
	Name() string
}

// CronQueueJob extends CronCommonJob to support queue-based job execution.
// Jobs implementing this interface can be triggered both by schedule (Run)
// and by queue messages (Work).
type CronQueueJob interface {
	CronCommonJob

	// Work processes a single job from the queue. This method is called when
	// a job is dequeued from Redis and needs to be executed.
	Work()
}

// Logger defines the logging interface required by the cronjob package.
// It supports structured logging at various levels.
type Logger interface {
	// Log writes a message at the specified level with the given arguments.
	Log(level logrus.Level, args ...interface{})

	// Logf writes a formatted message at the specified level.
	Logf(level logrus.Level, format string, args ...interface{})
}

// Alarm defines the interface for sending alerts when critical events occur.
// Implementations might send notifications to monitoring systems, messaging
// platforms, or other alerting mechanisms.
type Alarm interface {
	// Alarm sends an alert message at the specified severity level.
	// This is typically used for errors, timeouts, or other critical conditions.
	Alarm(level logrus.Level, message string)
}

// JobConfig holds the configuration needed by cron jobs for logging and alerting.
// It combines context, logger, and alarm interfaces to provide a complete
// execution environment for jobs.
type JobConfig struct {
	ctx    context.Context
	logger Logger
	alarm  Alarm
}

// NewJobConfig creates a new JobConfig with the provided context, logger, and alarm.
func NewJobConfig(ctx context.Context, logger Logger, alarm Alarm) *JobConfig {
	return &JobConfig{
		ctx:    ctx,
		logger: logger,
		alarm:  alarm,
	}
}

// Context returns the context associated with this job configuration.
func (jc *JobConfig) Context() context.Context {
	return jc.ctx
}

// Log writes a message at the specified level.
func (jc *JobConfig) Log(level logrus.Level, args ...interface{}) {
	jc.logger.Log(level, args...)
}

// Logf writes a formatted message at the specified level.
func (jc *JobConfig) Logf(level logrus.Level, format string, args ...interface{}) {
	jc.logger.Logf(level, format, args...)
}

// Alarm sends an alert at the specified level.
func (jc *JobConfig) Alarm(level logrus.Level, message string) {
	jc.alarm.Alarm(level, message)
}

// cronJob is an internal implementation of CronCommonJob and CronQueueJob.
// It wraps function pointers to provide a simple way to create jobs from
// existing functions.
type cronJob struct {
	run  func()
	work func()
	name string
}

// Run executes the run function.
func (f cronJob) Run() {
	f.run()
}

// Work executes the work function.
func (f cronJob) Work() {
	f.work()
}

// Name returns the job's unique identifier.
func (f cronJob) Name() string {
	return f.name
}
