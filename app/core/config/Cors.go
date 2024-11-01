package config

import (
	"encoding/json"
	"go_web/app/helper"
)

type Cors struct {
	Origin []string
}

func CorsConfigLoad() *Cors {
	var instance = &Cors{}
	err, str := helper.ReadJsonFile("config/cors.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, instance)
	return instance
}
