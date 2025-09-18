package queue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type QueueWorker struct {
	ctx         context.Context
	queueMap    map[string]*QueueJob
	queueRKey   []string
	redisClient *redis.Client
	queueLogger
	queueAlarm
}

func NewQueueWorker(ctx context.Context, redisClient *redis.Client, queueLogger queueLogger, queueAlarm queueAlarm) *QueueWorker {
	return &QueueWorker{
		ctx:         ctx,
		queueMap:    make(map[string]*QueueJob),
		queueRKey:   []string{},
		redisClient: redisClient,
		queueLogger: queueLogger,
		queueAlarm:  queueAlarm,
	}
}

func (q *QueueWorker) AddQueueJob(jobList []*QueueJob) (err error) {
	for _, job := range jobList {
		if _, ok := q.queueMap[job.Name]; ok {
			return fmt.Errorf("a task: [%s] already exists", job.Name)
		}

		for j, valLen := 0, q.redisClient.LLen(q.ctx, getPreConsumeRKey(job.Name)).Val(); j < int(valLen); j++ {
			rKeyVal, err := q.redisClient.RPopLPush(q.ctx, getPreConsumeRKey(job.Name), job.Name).Result()
			if err != nil {
				q.Alarm(logrus.ErrorLevel, fmt.Sprintf("Redis RPopLPush %s", job.Name))
				continue
			}
			q.Logf(logrus.InfoLevel, "redis insert %s %s", job.Name, rKeyVal)
		}
		q.queueRKey = append(q.queueRKey, job.Name)
		q.queueMap[job.Name] = job
		if job.Retry {
			q.queueRKey = append(q.queueRKey, getFailConsumeRKey(job.Name))
			q.queueMap[getFailConsumeRKey(job.Name)] = job
		}
	}

	return nil
}

func (q *QueueWorker) Start() {
	q.Log(logrus.InfoLevel, "start listen redis queue")

	// 用多个G监听
	go func() {
		cancelCtx, cancelFunc := context.WithCancel(q.ctx)
		defer func() {
			cancelFunc()
			if err := recover(); err != nil {
				q.Alarm(logrus.ErrorLevel, fmt.Sprintf("Listen Queue Error %v", err))
			}
		}()
		paramsChan := make(chan []string)
		// 启动work
		for i := 0; i < 20; i++ {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						q.Alarm(logrus.ErrorLevel, fmt.Sprintf("Work For Queue Error %v", err))
					}
				}()
				for {
					select {
					case param := <-paramsChan:
						if len(param) != 2 {
							continue
						}
						rKey := param[0]
						rKeyVal := param[1]
						err := q.consume(rKey, rKeyVal)
						if err != nil {
							q.Alarm(logrus.ErrorLevel, err.Error())
							continue
						}
					case <-cancelCtx.Done():
						return
					}
				}
			}()
		}

		for len(q.queueRKey) > 0 {
			val, err := q.listen()
			if err != nil {
				q.Alarm(logrus.ErrorLevel, fmt.Sprintf("Queue Error Redis BLPop Data Wrong %s", err.Error()))
				continue
			}
			if len(val) == 0 {
				continue
			}
			paramsChan <- val
		}
	}()
}

func (q *QueueWorker) listen() ([]string, error) {
	val, err := q.redisClient.BLPop(q.ctx, time.Second*20, q.queueRKey...).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("redis BLPop Fail: %s", err.Error())
	}
	// 根据返回的key 执行相应监听任务
	if len(val) != 2 {
		return nil, fmt.Errorf("redis Listen Data Wrong %+v", val)
	}
	return val, err
}

func (q *QueueWorker) consume(rKey, rKeyVal string) (err error) {
	defer func() {
		_, redisErr := q.redisClient.LRem(q.ctx, getPreConsumeRKey(rKey), 1, rKeyVal).Result()
		if redisErr != nil {
			err = fmt.Errorf("[%w] Redis LRem Fail %s %s", err, getPreConsumeRKey(rKey), rKeyVal)
		}
	}()
	// 插入预消费队列
	_, err = q.redisClient.RPush(q.ctx, getPreConsumeRKey(rKey), rKeyVal).Result()
	if err != nil {
		return fmt.Errorf("redis RPush %s %s", rKey, err.Error())
	}
	// 超时则提醒
	nexCtx, cancel := context.WithTimeout(q.ctx, time.Minute*1)
	defer cancel()
	go func() {
		<-nexCtx.Done()
		if errors.Is(nexCtx.Err(), context.DeadlineExceeded) {
			q.Alarm(logrus.InfoLevel, fmt.Sprintf("queue workForQueue [%s] [%s] work too slow", rKey, rKeyVal))
		}
	}()

	err = q.work(rKey, rKeyVal)
	if err != nil {
		// 若已重试 直接报错
		if errors.Is(errors.Unwrap(err), errRetry) || errors.Is(errors.Unwrap(err), errNotRetry) {
			return fmt.Errorf("work for queue error [%s] %s", rKey, err.Error())
		} else {
			// 发生错误 重新执行 从预消费队列推入重试队列
			_, err1 := q.redisClient.RPopLPush(q.ctx, getPreConsumeRKey(rKey), getFailConsumeRKey(rKey)).Result()
			if err1 != nil {
				return fmt.Errorf("[%w] Redis RPopLPush %s %s %s", err, rKey, rKeyVal, err1.Error())
			}
		}
	}
	return nil
}

func (q *QueueWorker) work(rKey, rKeyVal string) (err error) {
	job := q.queueMap[rKey]
	err = job.Run(rKeyVal)
	if err != nil {
		if rKey == getFailConsumeRKey(job.Name) {
			return fmt.Errorf("[%w] %w", errRetry, err)
		} else if job.Retry {
			return err
		} else {
			return fmt.Errorf("[%w] %w", errNotRetry, err)
		}
	}

	return nil
}
