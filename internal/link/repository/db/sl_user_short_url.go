package db

import (
	"context"
	"gorm.io/gorm"
	"short-link/database/mysql"
	"time"
)

type SlUserShortUrl struct {
	ID        uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ShortUrl  string `gorm:"column:short_url;NOT NULL"`            // 短链
	UserId    uint64 `gorm:"column:user_id;default:0;NOT NULL"`    // 用户ID
	CreatedAt int64  `gorm:"column:created_at;default:0;NOT NULL"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;default:0;NOT NULL"` // 更新时间
	IsDel     int    `gorm:"column:is_del;default:0;NOT NULL"`
}

type SlSlUserShortUrlDao struct {
	db *gorm.DB
}

func NewSlSlUserShortUrlDao(ctx context.Context, db ...*gorm.DB) *SlSlUserShortUrlDao {
	client := mysql.NewDBClient(ctx)
	if len(db) > 0 {
		client = db[0]
	}
	return &SlSlUserShortUrlDao{
		db: client,
	}
}

func (m *SlUserShortUrl) TableName() string {
	return "sl_user_short_url"
}

func (m *SlUserShortUrl) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().UnixMilli()
	m.UpdatedAt = time.Now().UnixMilli()
	return
}

func (m *SlSlUserShortUrlDao) Create(u *SlUserShortUrl, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Create(&u).Error
}

func (m *SlSlUserShortUrlDao) PageByUserId(userId uint64, lastId uint64, pageSize int) ([]*SlUserShortUrl, error) {
	var res []*SlUserShortUrl
	err := m.db.Table((&SlUserShortUrl{}).TableName()).
		Where("user_id = ? and id > ? and is_del = 0", userId, lastId).
		Limit(pageSize).Find(&res).Error
	return res, err
}

func (m *SlSlUserShortUrlDao) UpdateByShortUrl(shortUrl string, data map[string]interface{}, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlUserShortUrl{}).TableName()).Where("short_url = ?", shortUrl).Updates(data).Error
}

func (m *SlSlUserShortUrlDao) DeleteByShortUrl(shortUrl string, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlUserShortUrl{}).TableName()).Where("short_url = ?", shortUrl).Delete(&SlUserShortUrl{}).Error
}
