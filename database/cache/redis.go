package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"short-link/config"
)

var (
	redisMap = make(map[string]*redis.Client)
)

func InitializeRedisClient() {
	cfg := config.GetConfig()
	for key, redisCfg := range cfg.Database.Redis {
		rdb := NewRedis(redisCfg)
		redisMap[key] = rdb
	}
}

// TODO redis链接池配置
func NewRedis(redisCfg config.Redis) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	return rdb
}
