package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"short-link/database/cache"
	"short-link/internal/consts"
	"short-link/internal/link/domain"
	"short-link/internal/link/repository/db"
	"short-link/utils/netx"

	"gorm.io/gorm"
)

type IBlackListRepository interface {
	AddBlackList(ctx context.Context, shortURL string, IP uint32) error
	GetByID(ctx context.Context, id uint64) (*db.SlBlackList, error)
	Delete(ctx context.Context, id uint64) error
	PageBlackList(ctx context.Context, shortURL string, IP uint32, page, pageSize int) ([]*db.SlBlackList, int64, error)
	GetByShortUrl(ctx context.Context, shortURL string) ([]*db.SlBlackList, error)
	GetBlackListWithCache(ctx context.Context, shortUrl string) (domain.BlackLists, error)
	UpdateByID(ctx context.Context, id uint64, data map[string]any, tx ...*gorm.DB) error
	RefreshCache(ctx context.Context, shortURL string) error
}

type BlackListRepository struct {
}

func (b BlackListRepository) RefreshCache(ctx context.Context, shortURL string) error {
	redisKey := fmt.Sprintf(consts.RedisKeyShortURLBlackList, shortURL)
	rdb := cache.NewRedisTool(ctx)
	blackLists, err := b.GetByShortUrl(ctx, shortURL)
	if err != nil {
		return err
	}
	list := DBModelToDomain(blackLists)
	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	_, err = rdb.Set(ctx, redisKey, string(bytes))
	return err
}

func (b BlackListRepository) GetBlackListWithCache(ctx context.Context, shortUrl string) (domain.BlackLists, error) {
	var (
		blacklist domain.BlackLists
		rdb       = cache.NewRedisTool(ctx)
		redisKey  = fmt.Sprintf(consts.RedisKeyShortURLBlackList, shortUrl)
	)

	err := rdb.AutoFetch(ctx, redisKey, 0, &blacklist, func(ctx context.Context) (any, error) {
		bl, err := b.GetByShortUrl(ctx, shortUrl)
		if err != nil {
			return nil, err
		}
		list := DBModelToDomain(bl)
		return list, nil
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

func (b BlackListRepository) UpdateByID(ctx context.Context, id uint64, data map[string]any, tx ...*gorm.DB) error {
	return db.NewSlBlackListDao(ctx).UpdateByID(id, data, tx...)
}

func NewBlackListRepository() IBlackListRepository {
	return &BlackListRepository{}
}

func DBModelToDomain(models []*db.SlBlackList) domain.BlackLists {
	list := make(domain.BlackLists, 0, len(models))
	for _, item := range models {
		list = append(list, domain.BlackList{
			IP:     netx.IntToIP(item.IP),
			Status: item.Status,
		})
	}
	return list
}
