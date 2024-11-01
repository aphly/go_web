package helper

import (
	"errors"
	"os"
)

func FileIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 判断文件或者文件夹不存在
func IsFileDirExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ReadJsonFile(jsonFile string) (err error, jsonRes []byte) {
	isExist := FileIsExist(jsonFile)
	if !isExist {
		return errors.New("文件不存在"), nil
	}
	info, err := os.ReadFile(jsonFile)
	if err != nil {
		return errors.New("文件读取错误"), nil
	}
	return nil, info
}
