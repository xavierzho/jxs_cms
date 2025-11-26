/*
 *@Description:
	用于向外界发送告警信息
		小飞机
		...
*/

package message

import (
	"github.com/sirupsen/logrus"
)

const (
	TestId = iota
	CmsId  // 用于通知一般信息
	RiskId // 用于通知重要信息
)

type Alarm interface {
	// SendMsg 向指定频道发送消息
	SendMsg(msg string, msgType int32)
	// Alarm 打印日志并告警
	Alarm(level logrus.Level, msg string)
	// AlertErrorMsg 打印 Error 日志并告警
	AlertErrorMsg(msg string, msgType int32)
	// NotifyInfoMsg 打印 Info 日志并告警
	NotifyInfoMsg(msg string, msgType int32)
}
