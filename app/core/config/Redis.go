package config

import (
	"encoding/json"
	"go_web/app/helper"
)

type RedisGroup struct {
	Write []Redis
	Read  []Redis
}

type Redis struct {
	Addr     string
	Password string
	PoolSize int
	Retries  int
	Db       int
}

func RedisConfig() *map[string]RedisGroup {
	var instance = make(map[string]RedisGroup)
	err, str := helper.ReadJsonFile("config/redis.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, &instance)
	return &instance
}
