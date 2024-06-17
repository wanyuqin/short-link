package repository

import (
	"context"
	"gorm.io/gorm"
	"short-link/database/mysql"
	"short-link/internal/link/repository/db"
)

type ILinkRepository interface {
	AddLink(ctx context.Context, link *db.SlLink) error
	GetOriginal(ctx context.Context, original string) (*db.SlOriginalShortUrl, error)
	GetByShort(ctx context.Context, short string) (*db.SlOriginalShortUrl, error)
	GetByUserId(ctx context.Context, userId uint64, originUrl string, lastId uint64, pageSize int)
}

type LinkRepository struct {
}

func (l LinkRepository) GetByUserId(ctx context.Context, userId uint64, originUrl string, lastId uint64, pageSize int) {

	if originUrl == "" {
		shortUrlTable := make(map[string][]string)
		userShortUrls, err := db.NewSlSlUserShortUrlDao(ctx).PageByUserId(userId, lastId, pageSize)
		if err != nil {
			return
		}
		for _, shortUrl := range userShortUrls {
			if shortUrl.ShortUrl == "" {
				continue
			}
			slLink := db.SlLink{}
			shortUrlTable[slLink.TableName(shortUrl.ShortUrl)] = append(shortUrlTable[slLink.TableName(shortUrl.ShortUrl)], shortUrl.ShortUrl)
		}

	}
}

func (l LinkRepository) AddLink(ctx context.Context, link *db.SlLink) error {
	return mysql.NewDBClient(ctx).Transaction(func(tx *gorm.DB) error {
		err := db.NewSlLinkDao(ctx).Create(link, tx)
		if err != nil {
			return err
		}
		err = db.NewSlOriginalShortUrlDao(ctx).Create(&db.SlOriginalShortUrl{
			OriginalUrl: link.OriginUrl,
			ShortUrl:    link.ShortUrl,
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

func (l LinkRepository) GetOriginal(ctx context.Context, original string) (*db.SlOriginalShortUrl, error) {
	return db.NewSlOriginalShortUrlDao(ctx).GetByOriginalUrl(original)

}

func (l LinkRepository) GetByShort(ctx context.Context, short string) (*db.SlOriginalShortUrl, error) {
	return db.NewSlOriginalShortUrlDao(ctx).GetByShortUrl(short)
}

func NewLinkRepository() ILinkRepository {
	return &LinkRepository{}
}
