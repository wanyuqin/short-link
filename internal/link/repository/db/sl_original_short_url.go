package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"short-link/database/mysql"
)

type SlOriginalShortUrl struct {
	ID          int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ShortUrl    string `gorm:"column:short_url;NOT NULL"`            // 短链
	OriginalUrl string `gorm:"column:original_url;NOT NULL"`         // 短链
	CreatedAt   int64  `gorm:"column:created_at;default:0;NOT NULL"` // 创建时间
	UpdatedAt   int64  `gorm:"column:updated_at;default:0;NOT NULL"` // 更新时间
	IsDel       int    `gorm:"column:is_del;default:0;NOT NULL"`
}

func (m *SlOriginalShortUrl) TableName() string {
	return "sl_original_short_url"
}

type SlOriginalShortUrlDao struct {
	db *gorm.DB
}

func NewSlOriginalShortUrlDao(ctx context.Context, db ...*gorm.DB) *SlOriginalShortUrlDao {
	client := mysql.NewDBClient(ctx)
	if len(db) > 0 {
		client = db[0]
	}
	return &SlOriginalShortUrlDao{db: client}
}

func (m *SlOriginalShortUrlDao) GetByOriginalUrl(originalUrl string) (*SlOriginalShortUrl, error) {
	var res SlOriginalShortUrl
	err := m.db.Where("original_url = ?", originalUrl).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (m *SlOriginalShortUrlDao) GetByShortUrl(shortUrl string) (*SlOriginalShortUrl, error) {
	var res SlOriginalShortUrl
	err := m.db.Where("short_url = ?", shortUrl).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}
