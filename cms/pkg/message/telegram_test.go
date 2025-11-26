package message

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSendTelMsg(t *testing.T) {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	alarm := NewTelegramAlarm(l)
	TelegramSetting = &TelegramConfig{
		BotUrl:          "https://api.telegram.org/bot6200298283:AAEbyCDAv7YRanYCy1ss-8p0vLkO-5a6Uho/sendmessage",
		TestID:          "-898634515",
		CMSID:           "-948128710",
		RiskID:          "-952979102",
		TelSendMsgLevel: 40,
	}

	alarm.SendMsg("dev", RiskId)
	alarm.SendMsg("server", CmsId)

}
