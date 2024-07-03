package repository

import (
	"context"
	"fmt"
	"short-link/database/cache"
	"short-link/internal/consts"
	"short-link/internal/link/repository/db"
	"short-link/utils/netx"
)

type IBlackListRepository interface {
	AddBlackList(ctx context.Context, shortUrl string, ip uint32) error
	GetById(ctx context.Context, id uint64) (*db.SlBlackList, error)
	Delete(ctx context.Context, id uint64) error
	PageBlackList(ctx context.Context, shortUrl string, ip uint32, page, pageSize int) ([]*db.SlBlackList, int64, error)
	GetByShortUrl(ctx context.Context, shortUrl string) ([]*db.SlBlackList, error)
	GetBlackListWithCache(ctx context.Context, shortUrl string) ([]string, error)
}

type BlackListRepository struct {
}

func (b BlackListRepository) GetBlackListWithCache(ctx context.Context, shortUrl string) ([]string, error) {
	var (
		blacklist []string
		rdb       = cache.NewRedisTool(ctx)
		redisKey  = fmt.Sprintf(consts.RedisKeyShortUrlBlackList, shortUrl)
	)

	err := rdb.AutoFetch(ctx, redisKey, 0, &blacklist, func(ctx context.Context) (interface{}, error) {
		bl, err := b.GetByShortUrl(ctx, shortUrl)
		if err != nil {
			return nil, err
		}
		ips := make([]string, 0, len(bl))
		for _, item := range bl {
			ips = append(ips, netx.IntToIP(item.Ip))
		}
		return ips, nil
	})
	return blacklist, err
}

func (b BlackListRepository) AddBlackList(ctx context.Context, shortUrl string, ip uint32) error {
	return db.NewSlBlackListDao(ctx).Create(&db.SlBlackList{
		ShortUrl: shortUrl,
		Ip:       ip,
	})
}

func (b BlackListRepository) GetById(ctx context.Context, id uint64) (*db.SlBlackList, error) {
	return db.NewSlBlackListDao(ctx).GetById(id)
}

func (b BlackListRepository) Delete(ctx context.Context, id uint64) error {
	return db.NewSlBlackListDao(ctx).Delete(id)
}

func (b BlackListRepository) PageBlackList(ctx context.Context, shortUrl string, ip uint32, page, pageSize int) ([]*db.SlBlackList, int64, error) {
	return db.NewSlBlackListDao(ctx).List(shortUrl, ip, page, pageSize)
}

func (b BlackListRepository) GetByShortUrl(ctx context.Context, shortUrl string) ([]*db.SlBlackList, error) {
	return db.NewSlBlackListDao(ctx).GetByShortUrl(shortUrl)
}

func NewBlackListRepository() IBlackListRepository {
	return &BlackListRepository{}
}
