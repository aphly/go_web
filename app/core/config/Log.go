package config

import (
	"encoding/json"
	"go_web/app/helper"
)

type Log struct {
	Path          string
	MaxSize       int64 //m
	MaxBufferSize int64
}

func LogConfigLoad() *Log {
	var instance = &Log{}
	err, str := helper.ReadJsonFile("config/log.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, instance)
	return instance
}
