package connect

import (
	"github.com/redis/go-redis/v9"
	"go_web/app/core/config"
	"sync"
)

var Rediss = make(map[string]*redis.Client)
var RedisLocker sync.RWMutex

func Redis(config *config.Redis, key string) *redis.Client {
	RedisLocker.RLock()
	rd, ok := Rediss[key]
	if ok {
		RedisLocker.RUnlock()
		return rd
	}
	RedisLocker.RUnlock()

	RedisLocker.Lock()
	defer RedisLocker.Unlock()
	if _, ok1 := Rediss[key]; ok1 {
		return Rediss[key]
	}
	Rediss[key] = getRedis(config)
	return Rediss[key]
}

func getRedis(config *config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       config.Addr,
		Password:   config.Password,
		DB:         config.Db,
		PoolSize:   config.PoolSize,
		MaxRetries: config.Retries,
	})
}
