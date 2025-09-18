package local

import (
	"context"

	"data_backend/pkg/logger"
	"data_backend/pkg/queue"
)

var QueueWorker *queue.QueueWorker

// 初始化队列并启动
func SetupQueue() (err error) {
	ctx := context.WithValue(Ctx, logger.ModuleKey, Module.Add(".queue"))
	queueLog := Logger.WithContext(ctx)
	QueueWorker = queue.NewQueueWorker(ctx, RedisClient.Client, queueLog, NewAlarm(queueLog))

	return nil
}
