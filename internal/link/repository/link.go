package repository

import (
	"context"
	"short-link/internal/link/repository/db"
)

type ILinkRepository interface {
	AddLink(ctx context.Context, link *db.SlLink) error
	GetByOriginal(ctx context.Context, original string) (*db.SlLink, error)
	GetByShort(ctx context.Context, short string) (*db.SlLink, error)
}

type LinkRepository struct {
}

func (l LinkRepository) AddLink(ctx context.Context, link *db.SlLink) error {
	//TODO implement me
	panic("implement me")
}

func (l LinkRepository) GetByOriginal(ctx context.Context, original string) (*db.SlLink, error) {
	//TODO implement me
	panic("implement me")
}

func (l LinkRepository) GetByShort(ctx context.Context, short string) (*db.SlLink, error) {
	//TODO implement me
	panic("implement me")
}

func NewLinkRepository() ILinkRepository {
	return &LinkRepository{}
}
