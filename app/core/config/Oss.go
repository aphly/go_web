package config

import (
	"encoding/json"
	"go_web/app/helper"
)

// {
// 	"AccessKeyId":"",
// 	"AccessKeySecret":"",
// 	"Endpoint":"",
// 	"Bucket" : "",
// 	"IsCname":"false",
// 	"Url"  :""
// }

type Oss struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	Bucket          string
	IsCname         string
	Url             string
}

func OssConfigLoad() *Oss {
	var instance = Oss{}
	err, str := helper.ReadJsonFile("config/oss.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, &instance)
	return &instance
}
