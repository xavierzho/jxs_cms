package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"data_backend/pkg"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	FILE_LAYOUT         = "-2006-01-02"
	FILE_LAYOUT_PATTERN = `-[\d]{4}-[\d]{2}-[\d]{2}`
)

type fileHook struct {
	logger *lumberjack.Logger
	date   *time.Time
	lock   *sync.Mutex
	levels []logrus.Level
}

func NewFileHook(jackLog *lumberjack.Logger) (*fileHook, *fileHook) {
	if !checkLogFileName(jackLog.Filename) {
		panic(fmt.Sprintf("Invalid log file name: %s\n", jackLog.Filename))
	}

	dir := filepath.Dir(jackLog.Filename)
	filename := filepath.Base(jackLog.Filename)
	jackLogImportance := &lumberjack.Logger{
		Filename:  filepath.Join(dir, "importance", filename),
		MaxSize:   jackLog.MaxSize,
		LocalTime: jackLog.LocalTime,
	}
	now := time.Now().In(pkg.Location)

	hook := &fileHook{
		logger: jackLog,
		date:   &now,
		lock:   &sync.Mutex{},
		levels: logrus.AllLevels,
	}

	var importanceLevel []logrus.Level
	for _, item := range logrus.AllLevels {
		if item != logrus.InfoLevel && item != logrus.DebugLevel {
			importanceLevel = append(importanceLevel, item)
		}
	}
	importanceLogHook := &fileHook{
		logger: jackLogImportance,
		date:   &now,
		lock:   &sync.Mutex{},
		levels: importanceLevel,
	}

	return hook, importanceLogHook
}

func checkLogFileName(fileName string) bool {
	filename := filepath.Base(fileName)
	ext := filepath.Ext(filename)

	if len(filename)-len(ext) < len(FILE_LAYOUT) {
		return false
	}

	compile := regexp.MustCompile(FILE_LAYOUT_PATTERN)
	return compile.MatchString(filename[len(filename)-len(FILE_LAYOUT)-len(ext) : len(filename)-len(ext)])
}

func (hook *fileHook) Fire(entry *logrus.Entry) (err error) {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	err = hook.RotateLog()
	if err != nil {
		return err
	}
	// 打印日志到文件
	_, err = hook.logger.Write([]byte(line))
	if err != nil {
		return err
	}

	return nil
}

func (hook *fileHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *fileHook) RotateLog() (err error) {
	dir := filepath.Dir(hook.logger.Filename)
	filename := filepath.Base(hook.logger.Filename)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)-len(FILE_LAYOUT)]

	// 日期判断时间
	now := time.Now().In(pkg.Location)

	hook.lock.Lock()
	defer hook.lock.Unlock()
	// 分割日期
	if hook.date.Format(FILE_LAYOUT) != now.Format(FILE_LAYOUT) {
		// 修改日志文件 文件名。
		// 生成新的日志文件
		hook.logger.Filename = filepath.Join(dir, fmt.Sprintf("%s%s%s", prefix, now.Format(FILE_LAYOUT), ext))
		err = hook.logger.Rotate()
		if err != nil {
			hook.logger.Write([]byte(err.Error()))
			return err
		}
		hook.date = &now
	}
	return nil
}
