package repository

import (
	"context"
	"short-link/internal/link/repository/db"
)

type IBlackListRepository interface {
	AddBlackList(ctx context.Context, shortUrl string, ip uint32) error
	GetById(ctx context.Context, id uint64) (*db.SlBlackList, error)
	Delete(ctx context.Context, id uint64) error
	PageBlackList(ctx context.Context, shortUrl string, ip uint32, page, pageSize int) ([]*db.SlBlackList, int64, error)
}

type BlackListRepository struct {
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

func NewBlackListRepository() IBlackListRepository {
	return &BlackListRepository{}
}
