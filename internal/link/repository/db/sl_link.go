package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"short-link/internal/mysql"
	"time"
)

type SlLink struct {
	ID        int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	OriginUrl string `gorm:"column:origin_url;NOT NULL"`           // 原始链接
	ShortUrl  string `gorm:"column:short_url;NOT NULL"`            // 短链
	ExpiredAt int64  `gorm:"column:expired_at;default:0;NOT NULL"` // 过期时间
	CreatedAt int64  `gorm:"column:created_at;default:0;NOT NULL"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;default:0;NOT NULL"` // 更新时间
	IsDel     int    `gorm:"column:is_del;default:0;NOT NULL"`
}

func (m *SlLink) TableName() string {
	return "sl_link"
}

type SlLinkDao struct {
	db *gorm.DB
}

func NewSlLinkDao(ctx context.Context, db ...*gorm.DB) *SlLinkDao {
	client := mysql.NewClient(ctx)
	if len(db) > 0 {
		client = db[0]
	}
	return &SlLinkDao{db: client}
}

func (m *SlLink) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().UnixMilli()
	m.UpdatedAt = time.Now().UnixMilli()
	return
}

func (m *SlLinkDao) Create(l *SlLink) error {
	return m.db.Create(l).Error
}

func (m *SlLinkDao) GetByOriginUrl(originUrl string) (*SlLink, error) {
	var res SlLink
	err := m.db.Where("origin_url = ?", originUrl).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}
