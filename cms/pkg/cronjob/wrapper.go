package cronjob

import (
	"fmt"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// Wrapper is a function type that wraps a CronCommonJob with additional functionality.
// Wrappers can add cross-cutting concerns like logging, recovery, timeout detection, etc.
type Wrapper func(CronCommonJob) CronCommonJob

// CronChain represents a chain of job wrappers that are applied in sequence.
// Wrappers are applied in reverse order, so the last wrapper in the chain
// is the outermost wrapper (executes first).
type CronChain struct {
	wrappers []Wrapper
}

// NewCronChain creates a new CronChain from the provided wrappers.
//
// Example:
//
//	chain := NewCronChain(
//	  RecoverWrapper(config),         // Applied last (outermost)
//	  LoggerWrapper(config),          // Applied second
//	  SkipIfStillRunningWrapper(config), // Applied first (innermost)
//	)
func NewCronChain(wrappers ...Wrapper) CronChain {
	return CronChain{wrappers: wrappers}
}

// then applies all wrappers in the chain to the given job.
// Wrappers are applied in reverse order to maintain intuitive ordering.
func (c CronChain) then(j CronCommonJob) CronCommonJob {
	for i := len(c.wrappers) - 1; i >= 0; i-- {
		j = c.wrappers[i](j)
	}
	return j
}

// RecoverWrapper creates a wrapper that recovers from panics in job execution.
// When a panic occurs, it captures the stack trace and sends an alarm.
// This prevents a single job failure from crashing the entire scheduler.
//
// Usage:
//
//	chain := NewCronChain(RecoverWrapper(config))
func RecoverWrapper(config *JobConfig) Wrapper {
	return func(job CronCommonJob) CronCommonJob {
		return cronJob{
			run: recoverWrapper(job.Run, job.Name(), "run", config),
			work: wrapWorkIfQueueJob(job, func(workFn func()) func() {
				return recoverWrapper(workFn, job.Name(), "work", config)
			}),
			name: job.Name(),
		}
	}
}

// recoverWrapper wraps a function with panic recovery logic.
func recoverWrapper(fn func(), jobName, funcName string, config *JobConfig) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				const stackSize = 64 << 10 // 64KB stack buffer
				buf := make([]byte, stackSize)
				buf = buf[:runtime.Stack(buf, false)]

				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				errMsg := fmt.Sprintf("cron job execute %s %s panic: %v\n%s",
					jobName, funcName, err, string(buf))
				config.Alarm(logrus.ErrorLevel, errMsg)
			}
		}()
		fn()
	}
}

// LoggerWrapper creates a wrapper that logs job execution start and completion.
// Successful executions are logged at INFO level, while panics are re-thrown
// to be handled by RecoverWrapper.
//
// Usage:
//
//	chain := NewCronChain(
//	  RecoverWrapper(config),
//	  LoggerWrapper(config),  // Should be inside RecoverWrapper
//	)
func LoggerWrapper(config *JobConfig) Wrapper {
	return func(job CronCommonJob) CronCommonJob {
		return cronJob{
			run: loggerWrapper(job.Run, job.Name(), "run", config),
			work: wrapWorkIfQueueJob(job, func(workFn func()) func() {
				return loggerWrapper(workFn, job.Name(), "work", config)
			}),
			name: job.Name(),
		}
	}
}

// loggerWrapper wraps a function with execution logging.
func loggerWrapper(fn func(), jobName, funcName string, config *JobConfig) func() {
	return func() {
		config.Logf(logrus.InfoLevel, "cron job execute %s %s start", jobName, funcName)
		defer func() {
			if r := recover(); r != nil {
				panic(r) // Re-throw for RecoverWrapper to handle
			}
			config.Logf(logrus.InfoLevel, "cron job execute %s %s success", jobName, funcName)
		}()
		fn()
	}
}

// SkipIfStillRunningWrapper creates a wrapper that prevents concurrent execution
// of the same job. If a job is still running when its next scheduled execution
// occurs, the new execution is skipped.
//
// This is useful for long-running jobs where overlapping executions could
// cause resource conflicts or data inconsistencies.
//
// Usage:
//
//	chain := NewCronChain(SkipIfStillRunningWrapper(config))
func SkipIfStillRunningWrapper(config *JobConfig) Wrapper {
	return func(job CronCommonJob) CronCommonJob {
		return cronJob{
			run: skipIfStillRunningWrapper(job.Run, job.Name(), "run", config),
			work: wrapWorkIfQueueJob(job, func(workFn func()) func() {
				return skipIfStillRunningWrapper(workFn, job.Name(), "work", config)
			}),
			name: job.Name(),
		}
	}
}

// skipIfStillRunningWrapper wraps a function to prevent concurrent execution.
func skipIfStillRunningWrapper(fn func(), jobName, funcName string, config *JobConfig) func() {
	ch := make(chan struct{}, 1)
	ch <- struct{}{} // Initialize as available

	return func() {
		select {
		case token := <-ch:
			defer func() {
				ch <- token // Return token
				if r := recover(); r != nil {
					panic(r) // Re-throw panic
				}
			}()
			fn()
		default:
			config.Logf(logrus.InfoLevel, "cron job execute %s %s skip (still running)", jobName, funcName)
		}
	}
}

// TimeoutReminderWrapper creates a wrapper that sends an alarm if job execution
// exceeds the specified time limit. Note that this does NOT cancel the job;
// it only sends a notification.
//
// Parameters:
//   - config: Job configuration for logging and alarms
//   - timeLimit: Duration threshold for sending timeout alerts
//
// Usage:
//
//	chain := NewCronChain(
//	  TimeoutReminderWrapper(config, 5*time.Minute),
//	)
func TimeoutReminderWrapper(config *JobConfig, timeLimit time.Duration) Wrapper {
	return func(job CronCommonJob) CronCommonJob {
		return cronJob{
			run: timeoutReminderWrapper(job.Run, job.Name(), "run", config, timeLimit),
			work: wrapWorkIfQueueJob(job, func(workFn func()) func() {
				return timeoutReminderWrapper(workFn, job.Name(), "work", config, timeLimit)
			}),
			name: job.Name(),
		}
	}
}

// timeoutReminderWrapper wraps a function to measure execution time and alert on slow execution.
func timeoutReminderWrapper(fn func(), jobName, funcName string, config *JobConfig, timeLimit time.Duration) func() {
	return func() {
		startTime := time.Now()
		fn()
		elapsed := time.Since(startTime)

		if elapsed >= timeLimit {
			alarmMsg := fmt.Sprintf("cron job execute %s %s exceeded time limit: took %s (limit: %s)",
				jobName, funcName, elapsed, timeLimit)
			config.Alarm(logrus.ErrorLevel, alarmMsg)
		}
	}
}

// wrapWorkIfQueueJob is a helper function that conditionally wraps the Work method
// if the job implements CronQueueJob. This reduces code duplication in wrappers.
func wrapWorkIfQueueJob(job CronCommonJob, wrapper func(func()) func()) func() {
	if qJob, ok := job.(CronQueueJob); ok {
		return wrapper(qJob.Work)
	}
	return nil
}
