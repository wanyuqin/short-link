package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"short-link/config"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var (
	defaultDbName = "default"

	redisMap = make(map[string]*RedisTool)
)

type RedisTool struct {
	rdb *redis.Client
}

func InitializeRedisClient() {
	cfg := config.GetConfig()
	for key, redisCfg := range cfg.Database.Redis {
		rdb := NewRedis(redisCfg)
		redisMap[key] = rdb
	}
}

// TODO redis链接池配置
func NewRedis(redisCfg config.Redis) *RedisTool {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	tool := &RedisTool{
		rdb: rdb,
	}
	return tool
}

func NewRedisTool(ctx context.Context, key ...string) *RedisTool {
	dbName := defaultDbName
	if len(key) > 0 {
		dbName = key[0]
	}
	return redisMap[dbName]
}

type FetchFunc func(ctx context.Context) (any, error)

func (tool *RedisTool) AutoFetch(ctx context.Context, key string, ttl time.Duration, res any, fetchFun FetchFunc) error {
	data, err := tool.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 不存在
			source, err := tool.fetchFromSource(ctx, key, fetchFun)
			if err != nil && !errors.Is(err, redis.Nil) {
				return err
			}
			b, err := json.Marshal(source)
			if err != nil {
				return err
			}
			_, err = tool.rdb.Set(ctx, key, string(b), ttl).Result()
			if err != nil {
				return err
			}

			return json.Unmarshal(b, res)
		}
		return err
	}
	return json.Unmarshal([]byte(data), res)
}

func (tool *RedisTool) Del(ctx context.Context, key ...string) (int64, error) {
	return tool.rdb.Del(ctx, key...).Result()
}

func (tool *RedisTool) fetchFromSource(ctx context.Context, key string, fetch FetchFunc) (any, error) {
	g := singleflight.Group{}

	ret, err, _ := g.Do(key, func() (any, error) {
		go func() {
			time.Sleep(100 * time.Millisecond)
			g.Forget(key)
		}()
		return fetch(ctx)
	})

	return ret, err

}
