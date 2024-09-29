package config

import (
	"encoding/json"
	"go_web/app/helper"
)

type Http struct {
	Listen string
	Appkey []byte
}

func HttpConfigLoad() *Http {
	var instance = &Http{}
	err, str := helper.ReadJsonFile("config/http.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, instance)
	return instance
}
