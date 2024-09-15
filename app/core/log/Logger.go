package log

import (
	"bytes"
	"fmt"
	"go_web/app/core/config"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Loggers = make(map[string]*Logger)
var logLocker sync.RWMutex

type Logger struct {
	sync.Mutex
	Path         string // 保存路径 如 log
	Name         string // 日志项目名 如 mysql
	MaxSize      int64  // 单个日志文件的最大大小（MB）
	BufferWriter *BufferWriter
}

func (l *Logger) writelog(pre string, msg string) error {
	l.Lock()
	defer l.Unlock()
	dir, err := l.mkLogDir(pre)
	if err != nil {
		return err
	}
	path, err := l.lastLogFile(dir)
	if err != nil {
		return err
	}
	msg = time.Now().Format("15:04:05") + " " + msg
	err = l.AppendContent(path, msg)
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) mkLogDir(pre string) (string, error) {
	getwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := getwd + "/" + l.Path + "/" + l.Name + "/" + pre
	err = os.MkdirAll(dir, 0644)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func (l *Logger) lastLogFile(dir string) (string, error) {
	dateStr := time.Now().Format("2006-01-02")
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	today_max := 0
	for _, file := range files {
		if file.IsDir() {
		} else {
			base := l.BaseByFileName(file.Name())
			date_n := strings.Split(base, "_")
			if date_n[0] == dateStr {
				atoi, _ := strconv.Atoi(date_n[1])
				if atoi > today_max {
					today_max = atoi
				}
			}
		}
	}
	if today_max == 0 {
		return dir + "/" + dateStr + "_1.log", nil
	} else {
		oldPath := dir + "/" + dateStr + "_" + strconv.Itoa(today_max) + ".log"
		fileInfo, err := os.Stat(oldPath)
		if err != nil {
			return "", err
		}
		if fileInfo.Size() > l.MaxSize*1024*1024 {
			return dir + "/" + dateStr + "_" + strconv.Itoa(today_max+1) + ".log", nil
		}
		return oldPath, nil
	}
}

func (l *Logger) BaseByFileName(fileName string) string {
	for i := len(fileName) - 1; i >= 0 && fileName[i] != '/'; i-- {
		if fileName[i] == '.' {
			return fileName[:i]
		}
	}
	return ""
}

func (l *Logger) AppendContent(filePath string, msg any) error {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()
	_, err = logFile.WriteString(fmt.Sprintln(msg))
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) Info(msg any) error {
	err := l.writelog("/info", fmt.Sprint(msg))
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) Warn(msg any) error {
	err := l.writelog("/warn", fmt.Sprint(msg))
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) Error(msg any) error {
	err := l.writelog("/error", fmt.Sprint(msg))
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) Debug(msg any) error {
	err := l.writelog("/debug", fmt.Sprint(msg))
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) BW(pre string, msg any) error {
	dir, err := l.mkLogDir(pre)
	if err != nil {
		return err
	}
	filePath, err := l.lastLogFile(dir)
	if err != nil {
		return err
	}
	msgNew := time.Now().Format("15:04:05") + " " + fmt.Sprint(msg)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	l.BufferWriter.Writer = file
	_, err = l.BufferWriter.Write([]byte(msgNew))
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) Request(msg any) error {
	err := l.writelog("/request", fmt.Sprint(msg))
	if err != nil {
		return err
	}
	return nil
}

func NewLogger(config *config.Log, key string) *Logger {
	logLocker.RLock()
	logger, ok := Loggers[key]
	if ok {
		logLocker.RUnlock()
		return logger
	}
	logLocker.RUnlock()
	logLocker.Lock()
	defer logLocker.Unlock()

	if _, ok1 := Loggers[key]; ok1 {
		return Loggers[key]
	}
	Loggers[key] = getLogger(config, key)
	return Loggers[key]

}

func getLogger(config *config.Log, key string) *Logger {
	return &Logger{
		Path:    config.Path,
		Name:    key,
		MaxSize: config.MaxSize,
		BufferWriter: &BufferWriter{
			MaxBufferSize: config.MaxBufferSize,
			Buffer:        bytes.NewBuffer(make([]byte, 0, config.MaxBufferSize)),
		},
	}
}
