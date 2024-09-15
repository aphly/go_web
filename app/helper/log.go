package helper

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const MaxSize = 10 // m

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

func LogLastFile() (*os.File, error) {
	logDir, err := MkLogDir("")
	if err != nil {
		return nil, err
	}
	dateStr := time.Now().Format("2006-01-02")
	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		if fileInfo.Size() > MaxSize*1024*1024 {
			logPath = logDir + "/" + dateStr + "_" + strconv.Itoa(today_max+1) + ".log"
		}
		logPath = oldPath
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	defer logFile.Close()
	return logFile, nil
}

func BaseByFileName(fileName string) string {
	for i := len(fileName) - 1; i >= 0 && fileName[i] != '/'; i-- {
		if fileName[i] == '.' {
			return fileName[:i]
		}
	}
	return ""
}
