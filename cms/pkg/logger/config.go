package logger

type LoggerConfig struct {
	LogSavePath string
	LogFileName string
	LogFileExt  string
}

var LoggerSetting *LoggerConfig

// 用于日志中标识 调用的 module 路径
// should not use built-in type string as key for value; define your own type to avoid collisions
// 为避免 context.value 冲突 定义自定义类型(type CustomInfo string 会导致冲突)
type CustomInfo struct {
	value string
}

func NewCustomInfo(value string) CustomInfo {
	return CustomInfo{value}
}

func (c CustomInfo) String() string {
	return c.value
}

func (c CustomInfo) Add(str string) CustomInfo {
	c.value += str
	return c
}

// TODO 统一 customInfo 行为
var ModuleKey = CustomInfo{"module"}

const ( // 方便外部 Const 引用
	REQUEST_URL_KEY = "requestUrl"
	USER_ID_KEY     = "userID"
)
