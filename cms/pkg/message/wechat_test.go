package message

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestWeChat(t *testing.T) {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	alarm := NewWeChatAlarm(l)
	WeChatSetting = &WeChatConfig{
		WebHook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=267aed8c-af84-420e-9248-f53ed3cd7905",
	}

	alarm.SendMsg("test", 0)
}
