package repository

import (
	"context"
	"errors"
	"fmt"
	"short-link/database/cache"
	"short-link/database/mysql"
	"short-link/internal/consts"
	"short-link/internal/link/repository/db"
	"short-link/pkg/hot_key"
	"short-link/utils/gox"
	"sync"

	"gorm.io/gorm"
)

type ILinkRepository interface {
	AddLink(ctx context.Context, link *db.SlLink) error
	GetOriginal(ctx context.Context, original string, userId uint64) (*db.SlOriginalShortURL, error)
	GetByShort(ctx context.Context, short string) (*db.SlLink, error)
	PageUserLink(ctx context.Context, userId uint64, originUrl string, page int, pageSize int) ([]*db.SlLink, int64, error)
	UpdateByShort(ctx context.Context, shortUrl string, data map[string]any) error
	DeleteByShort(ctx context.Context, shortUrl string) error

	GetByShortWithCache(ctx context.Context, shortUrl string) (*db.SlLink, error)
}

type LinkRepository struct {
}

func (l LinkRepository) GetByShortWithCache(ctx context.Context, shortUrl string) (*db.SlLink, error) {
	var (
		redisKey = fmt.Sprintf(consts.RedisKeyShorURL, shortUrl)
		link     = db.SlLink{}
	)
	fetchFunc := func(ctx context.Context, key any) (any, error) {
		rdb := cache.NewRedisTool(ctx)
		res := db.SlLink{}
		err := rdb.AutoFetch(ctx, redisKey, 0, &res, func(ctx context.Context) (any, error) {
			short, err := l.GetByShort(ctx, shortUrl)
			if err != nil {
				return nil, err
			}
			if short == nil {
				return nil, errors.New("record not found")
			}
			return short, nil
		})
		return res, err
	}
	err := hot_key.GetHotKeyTool().Fetch(ctx, redisKey, &link, fetchFunc)

	return &link, err
}

func (l LinkRepository) DeleteByShort(ctx context.Context, shortUrl string) error {
	err := mysql.NewDBClient(ctx).Transaction(func(tx *gorm.DB) error {

		if err := db.NewSlLinkDao(ctx).DeleteByShort(shortUrl, tx); err != nil {
			return err
		}

		if err := db.NewSlSlUserShortURLDao(ctx).DeleteByShortURL(shortUrl, tx); err != nil {
			return err
		}

		if err := db.NewSlOriginalShortURLDao(ctx).DeleteByShortURL(shortUrl, tx); err != nil {
			return err
		}

		// 删除缓存
		rdb := cache.NewRedisTool(ctx)
		redisKey := fmt.Sprintf(consts.RedisKeyShorURL, shortUrl)
		if _, err := rdb.Del(ctx, redisKey); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (l LinkRepository) UpdateByShort(ctx context.Context, shortUrl string, data map[string]any) error {
	err := mysql.NewDBClient(ctx).Transaction(func(tx *gorm.DB) error {
		if err := db.NewSlLinkDao(ctx).UpdateByShortURL(shortUrl, data, tx); err != nil {
			return err
		}

		if err := db.NewSlSlUserShortURLDao(ctx).UpdateByShortURL(shortUrl, data, tx); err != nil {
			return err
		}

		if err := db.NewSlOriginalShortURLDao(ctx).UpdateByShortURL(shortUrl, data, tx); err != nil {
			return err
		}
		return nil

	})
	return err

}

func (l LinkRepository) PageUserLink(ctx context.Context, userId uint64, originUrl string, page int, pageSize int) ([]*db.SlLink, int64, error) {
	var (
		res   []*db.SlLink
		lock  = sync.Mutex{}
		total int64
	)
	urls, total, err := db.NewSlOriginalShortURLDao(ctx).PageByOriginalURL(originUrl, userId, page, pageSize)
	shortUrlTable := make(map[string][]string)
	if err != nil {
		return res, total, err
	}
	for _, url := range urls {
		if url.ShortURL == "" {
			continue
		}
		slLink := db.SlLink{}
		shortUrlTable[slLink.TableName(url.ShortURL)] = append(shortUrlTable[slLink.TableName(url.ShortURL)], url.ShortURL)
	}
	wg := gox.NewWaitGroup()
	for key, value := range shortUrlTable {
		urls := value
		table := key
		wg.RunSafe(ctx, func(ctx context.Context) {
			slLinks, err := db.NewSlLinkDao(ctx).GetByShortURLsWithTableName(table, urls)
			if err != nil {
				return
			}
			lock.Lock()
			defer lock.Unlock()
			res = append(res, slLinks...)

		})
	}
	wg.Wait()

	return res, total, nil
}

func (l LinkRepository) AddLink(ctx context.Context, link *db.SlLink) error {
	return mysql.NewDBClient(ctx).Transaction(func(tx *gorm.DB) error {
		err := db.NewSlLinkDao(ctx).Create(link, tx)
		if err != nil {
			return err
		}
		err = db.NewSlOriginalShortURLDao(ctx).Create(&db.SlOriginalShortURL{
			OriginURL: link.OriginURL,
			ShortURL:  link.ShortURL,
			UserID:    link.UserID,
		})
		if err != nil {
			return err
		}
		err = db.NewSlSlUserShortURLDao(ctx).Create(&db.SlUserShortURL{
			UserID:   link.UserID,
			ShortURL: link.ShortURL,
		})
		return err
	})
}

func (l LinkRepository) GetOriginal(ctx context.Context, original string, userId uint64) (*db.SlOriginalShortURL, error) {
	return db.NewSlOriginalShortURLDao(ctx).GetByOriginalURL(original, userId)

}

func (l LinkRepository) GetByShort(ctx context.Context, short string) (*db.SlLink, error) {
	return db.NewSlLinkDao(ctx).GetByShortURL(short)
}

func NewLinkRepository() ILinkRepository {
	return &LinkRepository{}
}
