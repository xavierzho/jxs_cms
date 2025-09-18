package queue

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type QueueJob struct {
	Name  string
	Retry bool
	Run   func(string) error
}

type queueLogger interface {
	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
}

type queueAlarm interface {
	Alarm(level logrus.Level, message string)
}

var errRetry = fmt.Errorf("retry Error")
var errNotRetry = fmt.Errorf("not retry Error")

// 预消费Key
func getPreConsumeRKey(key string) string {
	return fmt.Sprintf("Pre:%s", key)
}

// 重试消费Key
func getFailConsumeRKey(key string) string {
	return fmt.Sprintf("Fail:%s", key)
}
