package app

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go_web/app/core/config"
	"go_web/app/core/connect"
	"go_web/app/core/log"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var Config *config.Config

func Init() {
	Config = config.NewConfig()
}

func DbW(keys ...string) *gorm.DB {
	key := "default"
	if len(keys) > 0 {
		key = keys[0]
	}
	dbConfig := (*Config.Db)[key]
	if len(dbConfig.Write) <= 0 {
		panic("数据库写配置错误")
	}
	return connect.Mysql(&dbConfig.Write[0], key+"_write")
}

func DbR(keys ...string) *gorm.DB {
	key := "default"
	if len(keys) > 0 {
		key = keys[0]
	}
	dbConfig, ok := (*Config.Db)[key]
	if !ok {
		panic("数据库读配置错误:" + key)
	}
	dbLen := len(dbConfig.Read)
	if dbLen <= 0 {
		panic("数据库读配置错误")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randIndex := r.Intn(dbLen)
	return connect.Mysql(&dbConfig.Read[randIndex], fmt.Sprintf("%s_read_%s", key, randIndex))
}

// app.Log("sss").debug("wwwwww")
func Log(names ...string) *log.Logger {
	name := "default"
	if len(names) > 0 {
		name = names[0]
	}
	return log.NewLogger(Config.Log, name)
}

func RedisW(keys ...string) *redis.Client {
	key := "default"
	if len(keys) > 0 {
		key = keys[0]
	}
	RedisConfig := (*Config.Redis)[key]
	if len(RedisConfig.Write) <= 0 {
		panic("Redis写配置错误")
	}
	return connect.Redis(&RedisConfig.Write[0], key+"_write")
}

func RedisR(keys ...string) *redis.Client {
	key := "default"
	if len(keys) > 0 {
		key = keys[0]
	}
	RedisConfig, ok := (*Config.Redis)[key]
	if !ok {
		panic("Redis读配置错误:" + key)
	}
	dbLen := len(RedisConfig.Read)
	if dbLen <= 0 {
		panic("Redis读配置错误")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randIndex := r.Intn(dbLen)
	return connect.Redis(&RedisConfig.Read[randIndex], fmt.Sprintf("%s_read_%s", key, randIndex))
}
