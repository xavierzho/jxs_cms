/*
发送小飞机告警
*/
package message

import (
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type TelegramConfig struct {
	BotUrl          string
	TestID          string
	CMSID           string
	RiskID          string
	TelSendMsgLevel int32
}

var TelegramSetting *TelegramConfig

type telegramLogger interface {
	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
}

type TelegramAlarm struct {
	telegramLogger
}

func NewTelegramAlarm(log telegramLogger) *TelegramAlarm {
	return &TelegramAlarm{log}
}

func getMsgID(msgType int32) string {
	switch msgType {
	case CMS_ID:
		return TelegramSetting.CMSID
	case RISK_ID:
		return TelegramSetting.RiskID
	default:
		return TelegramSetting.CMSID
	}
}

func (t *TelegramAlarm) SendMsg(msg string, msgType int32) {
	if TelegramSetting.TelSendMsgLevel&(1<<(msgType-1)) != 0 {
		t.Logf(logrus.InfoLevel, "发送telegram信息, text: %s", msg)
		return
	}

	chatID := getMsgID(msgType)

	client := &http.Client{}
	//post请求
	postValues := url.Values{}
	postValues.Add("text", msg)
	postValues.Add("chat_id", chatID)

	resp, err := client.PostForm(TelegramSetting.BotUrl, postValues)
	if err != nil {
		t.Logf(logrus.ErrorLevel, "发送post请求失败: %s %s", msg, err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Logf(logrus.ErrorLevel, "发送telegram信息失败且解析返回失败: %s %v", msg, err.Error())
			return
		}
		t.Logf(logrus.ErrorLevel, "发送telegram信息失败: %s %v %s", msg, resp.StatusCode, string(b))
		return
	}

	t.Logf(logrus.InfoLevel, "发送telegram信息完成, text: %s", msg)
}

func (t *TelegramAlarm) Alarm(level logrus.Level, msg string) {
	switch level {
	case logrus.DebugLevel, logrus.InfoLevel:
		t.NotifyInfoMsg(msg, CMS_ID)
	case logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		t.AlertErrorMsg(msg, CMS_ID)
	}
}

func (t *TelegramAlarm) AlertErrorMsg(msg string, msgType int32) {
	t.Log(logrus.ErrorLevel, msg)
	t.SendMsg(msg, msgType)
}

func (t *TelegramAlarm) NotifyInfoMsg(msg string, msgType int32) {
	t.Log(logrus.InfoLevel, msg)
	t.SendMsg(msg, msgType)
}
