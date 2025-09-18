package message

type SMSConfig struct {
	BuKa struct {
		Url         string
		CallBackUrl string
		APPID       string
		CodeAPPID   string
		APIKey      string
		APISecret   string
	}
}

var SMSSetting *SMSConfig
