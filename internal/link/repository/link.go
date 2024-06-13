package repository

import (
	"context"
	"short-link/internal/link/repository/db"
)

type ILinkRepository interface {
	AddLink(ctx context.Context, link *db.SlLink) error
	GetOriginal(ctx context.Context, original string) (*db.SlOriginalShortUrl, error)
	GetByShort(ctx context.Context, short string) (*db.SlOriginalShortUrl, error)
}

type LinkRepository struct {
}

func (l LinkRepository) AddLink(ctx context.Context, link *db.SlLink) error {
	//TODO implement me
	panic("implement me")
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
