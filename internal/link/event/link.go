package event

import (
	"context"
	"fmt"
	"short-link/database/cache"
	"short-link/internal/consts"
)

func DeleteShortUrlEvent(ctx context.Context, shortUrl string) error {
	if shortUrl == "" {
		return fmt.Errorf("short url is empty")
	}
	rdb := cache.NewRedisTool(ctx)
	redisKey := fmt.Sprintf(consts.RedisKeyShortURLBlackList, shortUrl)
	_, err := rdb.Del(ctx, redisKey)
	return err
}
