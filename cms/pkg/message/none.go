/*
仅日志
*/
package message

import (
	"github.com/sirupsen/logrus"
)

type NoneConfig struct {
}

var NoneSetting *NoneConfig

type noneLogger interface {
	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
}

type NoneAlarm struct {
	noneLogger
}

func NewNoneAlarm(log noneLogger) *NoneAlarm {
	return &NoneAlarm{log}
}

func (t *NoneAlarm) SendMsg(msg string, msgType int32) {
}

func (t *NoneAlarm) Alarm(level logrus.Level, msg string) {
	switch level {
	case logrus.DebugLevel, logrus.InfoLevel:
		t.NotifyInfoMsg(msg, CMS_ID)
	case logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		t.AlertErrorMsg(msg, CMS_ID)
	}
}

func (t *NoneAlarm) AlertErrorMsg(msg string, msgType int32) {
	t.Log(logrus.ErrorLevel, msg)
	t.SendMsg(msg, msgType)
}

func (t *NoneAlarm) NotifyInfoMsg(msg string, msgType int32) {
	t.Log(logrus.InfoLevel, msg)
	t.SendMsg(msg, msgType)
}
