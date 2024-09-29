package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const MaxSize = 10 // m

var LogLock sync.Mutex

type WriterLog struct{}

func (this WriterLog) Printf(format string, v ...interface{}) {
	msg := time.Now().Format("2006-01-02 15:04:05") + " " + fmt.Sprintf(format, v...)
	LogLastFile(msg)
}

func MkLogDir(pre string) (string, error) {
	getwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := getwd + "/log/mysql/" + pre
	err = os.MkdirAll(dir, 0644)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func LogLastFile(msg any) error {
	LogLock.Lock()
	defer LogLock.Unlock()
	logDir, err := MkLogDir("")
	if err != nil {
		return err
	}
	dateStr := time.Now().Format("2006-01-02")
	files, err := os.ReadDir(logDir)
	if err != nil {
		return err
	}
	today_max := 0
	for _, file := range files {
		if file.IsDir() {
		} else {
			base := BaseByFileName(file.Name())
			date_n := strings.Split(base, "_")
			if date_n[0] == dateStr {
				atoi, _ := strconv.Atoi(date_n[1])
				if atoi > today_max {
					today_max = atoi
				}
			}
		}
	}
	logPath := ""
	if today_max == 0 {
		logPath = logDir + "/" + dateStr + "_1.log"
	} else {
		oldPath := logDir + "/" + dateStr + "_" + strconv.Itoa(today_max) + ".log"
		fileInfo, err := os.Stat(oldPath)
		if err != nil {
			return err
		}
		if fileInfo.Size() > MaxSize*1024*1024 {
			logPath = logDir + "/" + dateStr + "_" + strconv.Itoa(today_max+1) + ".log"
		}
		logPath = oldPath
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func BaseByFileName(fileName string) string {
	for i := len(fileName) - 1; i >= 0 && fileName[i] != '/'; i-- {
		if fileName[i] == '.' {
			return fileName[:i]
		}
	}
	return ""
}

func AppendContent(filePath string, msg any) error {
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
