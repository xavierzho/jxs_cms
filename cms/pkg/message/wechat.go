/*
发送企业微信告警
*/
package message

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type WeChatConfig struct {
	WebHook string
}

var WeChatSetting *WeChatConfig

type weChatLogger interface {
	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
}

type WeChatAlarm struct {
	weChatLogger
}

func NewWeChatAlarm(log weChatLogger) *WeChatAlarm {
	return &WeChatAlarm{log}
}

func (t *WeChatAlarm) SendMsg(msg string, msgType int32) {
	client := &http.Client{}
	//post请求
	msgMap := map[string]any{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg,
		},
	}

	msgBytes, err := json.Marshal(msgMap)
	if err != nil {
		t.Logf(logrus.ErrorLevel, "发送post请求时 json.Marshal 出错: %s %s", msg, err.Error())
		return
	}

	resp, err := client.Post(WeChatSetting.WebHook, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		t.Logf(logrus.ErrorLevel, "发送post请求失败: %s %s", msg, err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Logf(logrus.ErrorLevel, "发送weChat信息失败且解析返回失败: %s %v", msg, err.Error())
			return
		}
		t.Logf(logrus.ErrorLevel, "发送weChat信息失败: %s %v %s", msg, resp.StatusCode, string(b))
		return
	}

	t.Logf(logrus.InfoLevel, "发送weChat信息完成, text: %s", msg)
}

func (t *WeChatAlarm) Alarm(level logrus.Level, msg string) {
	switch level {
	case logrus.DebugLevel, logrus.InfoLevel:
		t.NotifyInfoMsg(msg, CMS_ID)
	case logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		t.AlertErrorMsg(msg, CMS_ID)
	}
}

func (t *WeChatAlarm) AlertErrorMsg(msg string, msgType int32) {
	t.Log(logrus.ErrorLevel, msg)
	t.SendMsg(msg, msgType)
}

func (t *WeChatAlarm) NotifyInfoMsg(msg string, msgType int32) {
	t.Log(logrus.InfoLevel, msg)
	t.SendMsg(msg, msgType)
}
