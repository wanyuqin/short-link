package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"short-link/database/cache"
	"short-link/database/mysql"
	"short-link/internal/consts"
	"short-link/internal/link/repository/db"
	"short-link/utils/gox"
	"sync"
)

type ILinkRepository interface {
	AddLink(ctx context.Context, link *db.SlLink) error
	GetOriginal(ctx context.Context, original string, userId uint64) (*db.SlOriginalShortUrl, error)
	GetByShort(ctx context.Context, short string) (*db.SlLink, error)
	PageUserLink(ctx context.Context, userId uint64, originUrl string, lastId uint64, pageSize int) ([]*db.SlLink, uint64, error)
	UpdateByShort(ctx context.Context, shortUrl string, data map[string]interface{}) error
	DeleteByShort(ctx context.Context, shortUrl string) error
}

type LinkRepository struct {
}

func (l LinkRepository) DeleteByShort(ctx context.Context, shortUrl string) error {

	err := mysql.NewDBClient(ctx).Transaction(func(tx *gorm.DB) error {

		if err := db.NewSlLinkDao(ctx).DeleteByShort(shortUrl, tx); err != nil {
			return err
		}

		if err := db.NewSlSlUserShortUrlDao(ctx).DeleteByShortUrl(shortUrl, tx); err != nil {
			return err
		}

		if err := db.NewSlOriginalShortUrlDao(ctx).DeleteByShortUrl(shortUrl, tx); err != nil {
			return err
		}

		// 删除缓存
		rdb := cache.NewRedisTool(ctx)
		redisKey := fmt.Sprintf(consts.RedisKeyShorUrl, shortUrl)
		if _, err := rdb.Del(ctx, redisKey); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (l LinkRepository) UpdateByShort(ctx context.Context, shortUrl string, data map[string]interface{}) error {
	err := mysql.NewDBClient(ctx).Transaction(func(tx *gorm.DB) error {
		if err := db.NewSlLinkDao(ctx).UpdateByShortUrl(shortUrl, data, tx); err != nil {
			return err
		}

		if err := db.NewSlSlUserShortUrlDao(ctx).UpdateByShortUrl(shortUrl, data, tx); err != nil {
			return err
		}

		if err := db.NewSlOriginalShortUrlDao(ctx).UpdateByShortUrl(shortUrl, data, tx); err != nil {
			return err
		}
		return nil

	})
	return err

}

func (l LinkRepository) PageUserLink(ctx context.Context, userId uint64, originUrl string, lastId uint64, pageSize int) ([]*db.SlLink, uint64, error) {
	var (
		res     []*db.SlLink
		lock    = sync.Mutex{}
		nLastId uint64
	)
	urls, err := db.NewSlOriginalShortUrlDao(ctx).PageByOriginalUrl(originUrl, userId, lastId, pageSize)
	shortUrlTable := make(map[string][]string)
	if err != nil {
		return res, nLastId, nil
	}
	for _, url := range urls {
		if url.ShortUrl == "" {
			continue
		}
		slLink := db.SlLink{}
		shortUrlTable[slLink.TableName(url.ShortUrl)] = append(shortUrlTable[slLink.TableName(url.ShortUrl)], url.ShortUrl)
		nLastId = url.ID
	}
	wg := &sync.WaitGroup{}
	for key, value := range shortUrlTable {
		urls := value
		table := key
		gox.RunSafe(ctx, wg, func(ctx context.Context) {
			slLinks, err := db.NewSlLinkDao(ctx).GetByShortUrlsWithTableName(table, urls)
			if err != nil {
				return
			}
			lock.Lock()
			lock.Unlock()
			res = append(res, slLinks...)

		})
	}
	wg.Wait()

	return res, nLastId, nil
}

func (l LinkRepository) AddLink(ctx context.Context, link *db.SlLink) error {
	return mysql.NewDBClient(ctx).Transaction(func(tx *gorm.DB) error {
		err := db.NewSlLinkDao(ctx).Create(link, tx)
		if err != nil {
			return err
		}
		err = db.NewSlOriginalShortUrlDao(ctx).Create(&db.SlOriginalShortUrl{
			OriginUrl: link.OriginUrl,
			ShortUrl:  link.ShortUrl,
			UserId:    link.UserId,
		})
		if err != nil {
			return err
		}
		err = db.NewSlSlUserShortUrlDao(ctx).Create(&db.SlUserShortUrl{
			UserId:   link.UserId,
			ShortUrl: link.ShortUrl,
		})
		return err
	})
}

func (l LinkRepository) GetOriginal(ctx context.Context, original string, userId uint64) (*db.SlOriginalShortUrl, error) {
	return db.NewSlOriginalShortUrlDao(ctx).GetByOriginalUrl(original, userId)

}

func (l LinkRepository) GetByShort(ctx context.Context, short string) (*db.SlLink, error) {
	return db.NewSlLinkDao(ctx).GetByShortUrl(short)
}

func NewLinkRepository() ILinkRepository {
	return &LinkRepository{}
}
