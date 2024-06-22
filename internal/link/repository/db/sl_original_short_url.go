package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"short-link/database/mysql"
	"time"
)

type SlOriginalShortUrl struct {
	ID        uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ShortUrl  string `gorm:"column:short_url;NOT NULL"`                             // 短链
	OriginUrl string `gorm:"column:origin_url;NOT NULL" json:"originUrl,omitempty"` // 原始链接
	UserId    uint64 `gorm:"column:user_id;NOT NULL"`                               // 用户ID
	CreatedAt int64  `gorm:"column:created_at;default:0;NOT NULL"`                  // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;default:0;NOT NULL"`                  // 更新时间
	IsDel     int    `gorm:"column:is_del;default:0;NOT NULL"`
}

func (m *SlOriginalShortUrl) TableName() string {
	return "sl_original_short_url"
}

type SlOriginalShortUrlDao struct {
	db *gorm.DB
}

func (m *SlOriginalShortUrl) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().UnixMilli()
	m.UpdatedAt = time.Now().UnixMilli()
	return
}

func NewSlOriginalShortUrlDao(ctx context.Context, db ...*gorm.DB) *SlOriginalShortUrlDao {
	client := mysql.NewDBClient(ctx)
	if len(db) > 0 {
		client = db[0]
	}
	return &SlOriginalShortUrlDao{db: client}
}

func (m *SlOriginalShortUrlDao) Create(u *SlOriginalShortUrl, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Create(&u).Error
}

func (m *SlOriginalShortUrlDao) GetByOriginalUrl(originalUrl string, userId uint64) (*SlOriginalShortUrl, error) {
	var res SlOriginalShortUrl
	err := m.db.Where("origin_url = ? and user_id = ? and  is_del = 0", originalUrl, userId).First(&res).Error
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
	err := m.db.Where("short_url = ? and is_del = 0", shortUrl).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (m *SlOriginalShortUrlDao) UpdateByShortUrl(shortUrl string, data map[string]interface{}, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlOriginalShortUrl{}).TableName()).Where("short_url = ?", shortUrl).Updates(data).Error
}

func (m *SlOriginalShortUrlDao) DeleteByShortUrl(shortUrl string, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlOriginalShortUrl{}).TableName()).Where("short_url = ?", shortUrl).Delete(&SlOriginalShortUrl{}).Error

}

func (m *SlOriginalShortUrlDao) PageByOriginalUrl(originalUrl string, userId uint64, lastId uint64, pageSize int) ([]*SlOriginalShortUrl, error) {
	var res []*SlOriginalShortUrl
	query := m.db.Where("user_id = ? ", userId)
	if originalUrl != "" {
		query = query.Where("origin_url = ? ", originalUrl)
	}
	if lastId > 0 {
		query = query.Where("id < ?", lastId)
	}
	err := query.Order("id desc").Limit(pageSize).Find(&res).Error
	return res, err
}
