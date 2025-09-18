package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"time"

	"data_backend/pkg"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 内嵌 logrus.Entry; 覆盖了 WithContext, WithField 方法
type Logger struct {
	*logrus.Entry
}

func NewLogger(ctx context.Context, jackLog *lumberjack.Logger) *Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(newFormatter())
	hook, importanceLogHook := NewFileHook(jackLog)
	l.AddHook(hook)
	l.AddHook(importanceLogHook)
	log := &Logger{logrus.NewEntry(l)}
	return log.WithContext(ctx)
}

func (l *Logger) AddHook(hook logrus.Hook) {
	l.Logger.AddHook(hook)
}

func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

// 调用 Entry.WithContext(ctx)
// 对于无请求的 *gin.Context 直接返回
// 其他情况则 WithField(ModuleKey) WithField(userID)
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.Entry = ll.Entry.WithContext(ctx)
	if v := ctx.Value(ModuleKey); v != nil {
		ll.WithField(ModuleKey.String(), v)
	}
	if v := ctx.Value(REQUEST_URL_KEY); v != nil {
		ll.WithField(REQUEST_URL_KEY, v)
	}
	if v := ctx.Value(USER_ID_KEY); v != nil {
		ll.WithField(USER_ID_KEY, v)
	}
	return ll
}

func (l *Logger) WithField(key string, value interface{}) {
	l.Entry = l.Entry.WithField(key, value)
}

func (l *Logger) WithCallersFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := make([]string, 0)
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		s := fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function)
		callers = append(callers, s)
		if !more {
			break
		}
	}
	ll := l.clone()
	ll.WithField("callers", callers)
	return ll
}

func (l *Logger) Value(key any) any {
	return l.Context.Value(key)
}

func (l *Logger) ModuleKey() CustomInfo {
	key := l.Value(ModuleKey)
	if key == nil {
		return CustomInfo{""}
	}

	return key.(CustomInfo)
}

type loggerFormatter struct {
	logrus.TextFormatter
}

var formatSortKeyMap = map[string]int{
	logrus.FieldKeyTime:        -8,
	logrus.FieldKeyLevel:       -7,
	REQUEST_URL_KEY:            -6,
	ModuleKey.String():         -5,
	logrus.FieldKeyMsg:         -4,
	logrus.FieldKeyLogrusError: -3,
	logrus.FieldKeyFile:        -2,
	logrus.FieldKeyFunc:        -1,
}

func newFormatter() *loggerFormatter {
	format := new(loggerFormatter)
	format.DisableTimestamp = true
	format.SortingFunc = func(keys []string) {
		sort.Slice(keys, func(i, j int) bool {
			if formatSortKeyMap[keys[i]] != formatSortKeyMap[keys[j]] {
				return formatSortKeyMap[keys[i]] < formatSortKeyMap[keys[j]]
			} else {
				return keys[i] < keys[j]
			}
		})
	}

	return format
}

func (f *loggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b, err := f.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	strByte := []byte(fmt.Sprintf("[%s] ", entry.Time.In(pkg.Location).Format(time.RFC3339)))
	strByte = append(strByte, b...)
	return strByte, err
}
