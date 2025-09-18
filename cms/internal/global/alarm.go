package global

import (
	"data_backend/pkg/logger"
	"data_backend/pkg/message"
)

var Alarm message.Alarm

func SetupAlarm() (err error) {
	Alarm = NewAlarm(Logger)

	return nil
}

func NewAlarm(log *logger.Logger) message.Alarm {
	return message.NewWeChatAlarm(log)
}
