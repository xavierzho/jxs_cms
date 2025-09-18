package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestCheckName(t *testing.T) {
	fmt.Println(checkLogFileName("test-2006-01-02"))
}

// 测试日志多路径输出
func TestOutput(t *testing.T) {
	jackLog := &lumberjack.Logger{
		Filename:  "test-0000-00-00.log",
		MaxSize:   500,
		LocalTime: true,
	}
	log := NewLogger(context.Background(), jackLog)
	log.Info("info 1")
	log.Error("Err 1")
	log.Info("info 2")
}

func TestColor(t *testing.T) {
	jackLog := &lumberjack.Logger{
		Filename:  "test-0000-00-00.log",
		MaxSize:   500,
		LocalTime: true,
	}
	l := NewLogger(context.Background(), jackLog)
	l.Info("\033[32mTest\033[0m\033[32m[info]\033[0m")
	ll := log.New(os.Stdout, "\r\n", log.LstdFlags)
	ll.Printf("\033[32mTest\033[0m\033[32m[info]\033[0m")
	file, _ := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND, os.ModePerm)
	defer file.Close()
	ll.SetOutput(file)
	ll.Printf("\033[32mTest\033[0m\033[32m[info]\033[0m")
}

func TestFormatSortFunc(t *testing.T) {
	format := newFormatter()
	keys := []string{
		logrus.FieldKeyLevel,
		logrus.FieldKeyLogrusError,
		logrus.FieldKeyTime,
		logrus.FieldKeyMsg,
		ModuleKey.String(),
		"a",
		logrus.FieldKeyFunc,
		"b",
		logrus.FieldKeyFile,
		"c",
	}
	format.SortingFunc(keys)
	fmt.Println(keys)
}
