package event

import (
	"context"
	"fmt"
	"short-link/database/cache"
	"short-link/internal/consts"
)

type UpdateBlackListStatusEventMsg struct {
	ShortURL string `json:"shortUrl"`
	Ip       string `json:"ip"`
	Status   int    `json:"status"`
}

// UpdateBlackListStatusEvent 黑名单改变
func UpdateBlackListStatusEvent(ctx context.Context, msg UpdateBlackListStatusEventMsg) error {
	if msg.ShortURL == "" || msg.Ip == "" {
		return nil
	}

	if msg.Status == consts.IPStatusActive {
		// 启用加入redis
		redisKey := fmt.Sprintf(consts.RedisKeyShortURLBlackList, msg.ShortURL)
		rdb := cache.NewRedisTool(ctx)
		_, err := rdb.GetClient().SAdd(ctx, redisKey, msg.Ip).Result()
		return err
	}

	if msg.Status == consts.IPStatusDisable {
		return deleteIp(ctx, msg.ShortURL, msg.Ip)
	}
	return nil
}

type DeleteBlackListEventMsg struct {
	ShortURL string `json:"shortUrl"`
	Ip       string `json:"ip"`
}

// DeleteBlackListEvent 黑名单删除
func DeleteBlackListEvent(ctx context.Context, msg DeleteBlackListEventMsg) error {
	if msg.ShortURL == "" || msg.Ip == "" {
		return nil
	}
	return deleteIp(ctx, msg.ShortURL, msg.Ip)
}

func deleteIp(ctx context.Context, shortURL string, ip string) error {
	redisKey := fmt.Sprintf(consts.RedisKeyShortURLBlackList, shortURL)
	rdb := cache.NewRedisTool(ctx)
	_, err := rdb.GetClient().SRem(ctx, redisKey, ip).Result()
	return err

}
