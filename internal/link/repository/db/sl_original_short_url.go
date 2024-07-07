package db

import (
	"context"
	"errors"
	"short-link/database/mysql"
	"time"

	"gorm.io/gorm"
)

type SlOriginalShortURL struct {
	ID        uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ShortURL  string `gorm:"column:short_url;NOT NULL"`                             // 短链
	OriginURL string `gorm:"column:origin_url;NOT NULL" json:"originUrl,omitempty"` // 原始链接
	UserID    uint64 `gorm:"column:user_id;NOT NULL"`                               // 用户ID
	CreatedAt int64  `gorm:"column:created_at;default:0;NOT NULL"`                  // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;default:0;NOT NULL"`                  // 更新时间
	IsDel     int    `gorm:"column:is_del;default:0;NOT NULL"`
}

func (_ *SlOriginalShortURL) TableName() string {
	return "sl_original_short_url"
}

type SlOriginalShortURLDao struct {
	db *gorm.DB
}

func (m *SlOriginalShortURL) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().UnixMilli()
	m.UpdatedAt = time.Now().UnixMilli()
	return
}

func NewSlOriginalShortURLDao(ctx context.Context, db ...*gorm.DB) *SlOriginalShortURLDao {
	client := mysql.NewDBClient(ctx)
	if len(db) > 0 {
		client = db[0]
	}
	return &SlOriginalShortURLDao{db: client}
}

func (m *SlOriginalShortURLDao) Create(u *SlOriginalShortURL, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Create(&u).Error
}

func (m *SlOriginalShortURLDao) GetByOriginalURL(originalURL string, userID uint64) (*SlOriginalShortURL, error) {
	var res SlOriginalShortURL
	err := m.db.Where("origin_url = ? and user_id = ? and  is_del = 0", originalURL, userID).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (m *SlOriginalShortURLDao) GetByShortURL(ShortURL string) (*SlOriginalShortURL, error) {
	var res SlOriginalShortURL
	err := m.db.Where("short_url = ? and is_del = 0", ShortURL).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (m *SlOriginalShortURLDao) UpdateByShortURL(ShortURL string, data map[string]any, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlOriginalShortURL{}).TableName()).Where("short_url = ?", ShortURL).Updates(data).Error
}

func (m *SlOriginalShortURLDao) DeleteByShortURL(ShortURL string, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlOriginalShortURL{}).TableName()).Where("short_url = ?", ShortURL).Delete(&SlOriginalShortURL{}).Error
}

func (m *SlOriginalShortURLDao) PageByOriginalURL(originalURL string, userId uint64, page int, pageSize int) ([]*SlOriginalShortURL, int64, error) {
	var (
		res   []*SlOriginalShortURL
		count int64
	)
	query := m.db.Table((&SlOriginalShortURL{}).TableName()).
		Where("user_id = ? ", userId)
	if originalURL != "" {
		query = query.Where("origin_url = ? ", originalURL)
	}
	err := query.Order("id desc").
		Count(&count).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&res).Error
	return res, count, err
}
