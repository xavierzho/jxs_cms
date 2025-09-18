/*
为适配运行、日志需求封装"github.com/robfig/cron/v3"
1.在原cron.job基础上增加了两个job接口
2.修改日志输出形式（对接global.Logger）
3.增加新的装饰器
*/
package cronjob

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CronCommonJob interface {
	Run()
	Name() string
}

type CronQueueJob interface {
	CronCommonJob
	Work()
}

type cronJob struct {
	run  func()
	work func()
	name string
}

func (f cronJob) Run() {
	f.run()
}
func (f cronJob) Work() {
	f.work()
}
func (f cronJob) Name() string {
	return f.name
}

type cronLogger interface {
	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
}

type cronAlarm interface {
	Alarm(level logrus.Level, message string)
}

type cronMessage struct {
	ctx context.Context
	cronLogger
	cronAlarm
}

func NewCronMessage(ctx context.Context, log cronLogger, alarm cronAlarm) *cronMessage {
	return &cronMessage{
		ctx,
		log,
		alarm,
	}
}
