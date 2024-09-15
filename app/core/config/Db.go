package config

import (
	"encoding/json"
	"go_web/app/helper"
)

type DbGroup struct {
	Write []Db
	Read  []Db
}

type Db struct {
	Host           string
	Port           int
	Database       string
	Username       string
	Password       string
	Charset        string
	TimeOut        int
	WriteTimeOut   int
	ReadTimeOut    int
	MaxIdleConnect int
	MaxOpenConnect int
}

func DbConfig() *map[string]DbGroup {
	var instance = make(map[string]DbGroup)
	err, str := helper.ReadJsonFile("config/db.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, &instance)
	return &instance
}
