package cronjob

import (
	"fmt"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

type CronJobWrapper func(CronCommonJob) CronCommonJob

type CronChain struct {
	wrappers []CronJobWrapper
}

func NewCronChain(c ...CronJobWrapper) CronChain {
	return CronChain{c}
}

func (c CronChain) then(j CronCommonJob) CronCommonJob {
	for i := range c.wrappers {
		j = c.wrappers[len(c.wrappers)-1-i](j)
	}
	return j
}

func RecoverWrapper(message *cronMessage) CronJobWrapper {
	return func(job CronCommonJob) CronCommonJob {
		funcJob := cronJob{
			run:  recoverWrapper(job.Run, job.Name(), "run", message),
			name: job.Name(),
		}
		if _, ok := job.(CronQueueJob); ok {
			funcJob.work = recoverWrapper(job.(CronQueueJob).Work, job.Name(), "work", message)
		}
		return funcJob
	}
}

func recoverWrapper(f func(), jobName, funName string, message *cronMessage) func() {
	return func() {
		// recover
		defer func() {
			if r := recover(); r != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				message.Alarm(logrus.ErrorLevel, fmt.Sprintf("cron job execute %s %s occur error: %v %s", jobName, funName, err, string(buf)))
			}
		}()
		f()
	}
}

func LoggerWrapper(message *cronMessage) CronJobWrapper {
	return func(job CronCommonJob) CronCommonJob {
		funcJob := cronJob{
			run:  loggerWrapper(job.Run, job.Name(), "run", message),
			name: job.Name(),
		}
		if _, ok := job.(CronQueueJob); ok {
			funcJob.work = loggerWrapper(job.(CronQueueJob).Work, job.Name(), "work", message)
		}
		return funcJob
	}
}

func loggerWrapper(f func(), jobName, funName string, message *cronMessage) func() {
	return func() {
		message.Logf(logrus.InfoLevel, "cron job execute %s %s start", jobName, funName)
		defer func() {
			if r := recover(); r != nil {
				panic(r) // 不做处理直接抛出
			} else {
				message.Logf(logrus.InfoLevel, "cron job execute %s %s success", jobName, funName)
			}
		}()
		f()
	}
}

func SkipIfStillRunningWrapper(message *cronMessage) CronJobWrapper {
	return func(job CronCommonJob) CronCommonJob {
		funcJob := cronJob{
			run:  skipIfStillRunningWrapper(job.Run, job.Name(), "run", message),
			name: job.Name(),
		}
		if _, ok := job.(CronQueueJob); ok {
			funcJob.work = skipIfStillRunningWrapper(job.(CronQueueJob).Work, job.Name(), "work", message)
		}
		return funcJob
	}
}

func skipIfStillRunningWrapper(f func(), jobName, funName string, message *cronMessage) func() {
	var ch = make(chan struct{}, 1)
	ch <- struct{}{}
	return func() {
		select {
		case v := <-ch:
			defer func() {
				ch <- v
				if r := recover(); r != nil {
					panic(r) // 继续抛出
				}
			}()
			f()
		default:
			message.Logf(logrus.InfoLevel, "cron job execute %s %s skip", jobName, funName)
		}
	}
}

func TimeoutReminderWrapper(message *cronMessage, timeLimit time.Duration) CronJobWrapper {
	return func(job CronCommonJob) CronCommonJob {
		funcJob := cronJob{
			run:  timeoutReminderWrapper(job.Run, job.Name(), "run", message, timeLimit),
			name: job.Name(),
		}
		if _, ok := job.(CronQueueJob); ok {
			funcJob.work = timeoutReminderWrapper(job.(CronQueueJob).Work, job.Name(), "work", message, timeLimit)
		}
		return funcJob
	}
}

func timeoutReminderWrapper(f func(), jobName, funName string, message *cronMessage, timeLimit time.Duration) func() {
	return func() {
		startTime := time.Now()
		f()
		subTime := time.Since(startTime)
		if subTime >= timeLimit {
			message.Alarm(logrus.ErrorLevel, fmt.Sprintf("cron job execute %s %s too slow: %s", jobName, funName, subTime.String()))
		}

	}
}
