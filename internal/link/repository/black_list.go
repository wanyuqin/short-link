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
	AddBlackList(ctx context.Context, shortURL string, IP uint32) error
	GetByID(ctx context.Context, id uint64) (*db.SlBlackList, error)
	Delete(ctx context.Context, id uint64) error
	PageBlackList(ctx context.Context, shortURL string, IP uint32, page, pageSize int) ([]*db.SlBlackList, int64, error)
	GetByShortUrl(ctx context.Context, shortURL string) ([]*db.SlBlackList, error)
	GetBlackListWithCache(ctx context.Context, shortUrl string) ([]string, error)
}

type BlackListRepository struct {
}

func (b BlackListRepository) GetBlackListWithCache(ctx context.Context, shortUrl string) ([]string, error) {
	var (
		blacklist []string
		rdb       = cache.NewRedisTool(ctx)
		redisKey  = fmt.Sprintf(consts.RedisKeyShortURLBlackList, shortUrl)
	)

	err := rdb.AutoFetch(ctx, redisKey, 0, &blacklist, func(ctx context.Context) (any, error) {
		bl, err := b.GetByShortUrl(ctx, shortUrl)
		if err != nil {
			return nil, err
		}
		IPs := make([]string, 0, len(bl))
		for _, item := range bl {
			IPs = append(IPs, netx.IntToIP(item.IP))
		}
		return IPs, nil
	})
	return blacklist, err
}

func (b BlackListRepository) AddBlackList(ctx context.Context, shortURL string, ip uint32) error {
	return db.NewSlBlackListDao(ctx).Create(&db.SlBlackList{
		ShortURL: shortURL,
		IP:       ip,
	})
}

func (b BlackListRepository) GetByID(ctx context.Context, id uint64) (*db.SlBlackList, error) {
	return db.NewSlBlackListDao(ctx).GetByID(id)
}

func (b BlackListRepository) Delete(ctx context.Context, id uint64) error {
	return db.NewSlBlackListDao(ctx).Delete(id)
}

func (b BlackListRepository) PageBlackList(ctx context.Context, shortUrl string, ip uint32, page, pageSize int) ([]*db.SlBlackList, int64, error) {
	return db.NewSlBlackListDao(ctx).List(shortUrl, ip, page, pageSize)
}

func (b BlackListRepository) GetByShortUrl(ctx context.Context, shortUrl string) ([]*db.SlBlackList, error) {
	return db.NewSlBlackListDao(ctx).GetByShortURL(shortUrl)
}

func NewBlackListRepository() IBlackListRepository {
	return &BlackListRepository{}
}
