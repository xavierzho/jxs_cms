package cronjob

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type JobWorker struct {
	ctx               context.Context
	queueRKey         string
	cronMessage       *cronMessage
	redisClient       *redis.Client
	jobList           map[string][]CronCommonJob
	frequentlyJobList map[string][]CronCommonJob
	queueJobDict      map[string]CronQueueJob
}

func NewJobWorker(ctx context.Context, queueRKey string, cronMessage *cronMessage, redisClient *redis.Client) *JobWorker {
	return &JobWorker{
		ctx:               ctx,
		queueRKey:         queueRKey,
		cronMessage:       cronMessage,
		redisClient:       redisClient,
		jobList:           make(map[string][]CronCommonJob),
		frequentlyJobList: make(map[string][]CronCommonJob),
		queueJobDict:      make(map[string]CronQueueJob),
	}
}

func (j *JobWorker) QueueRKey() string {
	return j.queueRKey
}

func (j *JobWorker) addQueueJobDict(job CronQueueJob) error {
	if _, ok := j.queueJobDict[job.Name()]; ok {
		return fmt.Errorf("a job: [%s] already exists", job.Name())
	}
	j.queueJobDict[job.Name()] = job

	return nil
}

func (j *JobWorker) AddJobList(jobList map[string][]CronCommonJob) error {
	for spec, jobs := range jobList {
		j.jobList[spec] = append(j.jobList[spec], jobs...)
	}

	return nil
}

func (j *JobWorker) StartJob(cronChain, queueCronChain CronChain) error {
	if len(j.jobList) == 0 {
		return nil
	}
	c := cron.New()
	return j.addJob(j.jobList, c, cronChain, queueCronChain)
}

func (j *JobWorker) AddFrequentlyJobList(jobList map[string][]CronCommonJob) error {
	for spec, jobs := range jobList {
		j.frequentlyJobList[spec] = append(j.frequentlyJobList[spec], jobs...)
	}

	return nil
}

func (j *JobWorker) StartFrequentlyJob(cronChain CronChain) error {
	if len(j.frequentlyJobList) == 0 {
		return nil
	}
	c := cron.New(cron.WithSeconds())
	return j.addJob(j.frequentlyJobList, c, cronChain, cronChain)
}

func (j *JobWorker) addJob(jobList map[string][]CronCommonJob, c *cron.Cron, cronChain CronChain, queueCronChain CronChain) (err error) {
	for spec, jobs := range jobList {
		for i := 0; i < len(jobs); i++ {
			var job CronCommonJob
			if _, ok := jobs[i].(CronQueueJob); ok {
				job = queueCronChain.then(jobs[i])
				j.addQueueJobDict(job.(CronQueueJob))
			} else {
				job = cronChain.then(jobs[i])
			}
			_, err = c.AddJob(spec, job)
			if err != nil {
				c.Stop()
				return fmt.Errorf("cron job execute %s AddJob Fail: %s", job.Name(), err.Error())
			}
		}
	}
	c.Start()
	return nil
}

// 用于向队列添加值 触发任务执行
func (j *JobWorker) AddJobToQueue(jobName string) {
	_, err := j.redisClient.RPush(j.ctx, j.queueRKey, jobName).Result()
	if err != nil {
		errMsg := fmt.Sprintf("cron job execute %s addJobToQueue Fail: %s", jobName, err.Error())
		j.cronMessage.Alarm(logrus.ErrorLevel, errMsg)
		return
	}
}

func (j *JobWorker) QueueJobRun(jobName string) error {
	if job, ok := j.queueJobDict[jobName]; ok {
		job.Work()
	} else {
		j.cronMessage.Logf(logrus.ErrorLevel, "cron job execute %s work Job not exit", jobName)
	}
	return nil
}
