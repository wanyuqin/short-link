package db

import (
	"context"
	"short-link/database/mysql"
	"time"

	"gorm.io/gorm"
)

type SlUserShortURL struct {
	ID        uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ShortURL  string `gorm:"column:short_url;NOT NULL"`            // 短链
	UserID    uint64 `gorm:"column:user_id;default:0;NOT NULL"`    // 用户ID
	CreatedAt int64  `gorm:"column:created_at;default:0;NOT NULL"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;default:0;NOT NULL"` // 更新时间
	IsDel     int    `gorm:"column:is_del;default:0;NOT NULL"`
}

type SlSlUserShortURLDao struct {
	db *gorm.DB
}

func NewSlSlUserShortURLDao(ctx context.Context, db ...*gorm.DB) *SlSlUserShortURLDao {
	client := mysql.NewDBClient(ctx)
	if len(db) > 0 {
		client = db[0]
	}
	return &SlSlUserShortURLDao{
		db: client,
	}
}

func (m *SlUserShortURL) TableName() string {
	return "sl_user_short_url"
}

func (m *SlUserShortURL) BeforeCreate(_ *gorm.DB) (err error) {
	m.CreatedAt = time.Now().UnixMilli()
	m.UpdatedAt = time.Now().UnixMilli()
	return
}

func (m *SlSlUserShortURLDao) Create(u *SlUserShortURL, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Create(&u).Error
}

func (m *SlSlUserShortURLDao) PageByUserId(userId uint64, lastId uint64, pageSize int) ([]*SlUserShortURL, error) {
	var res []*SlUserShortURL
	err := m.db.Table((&SlUserShortURL{}).TableName()).
		Where("user_id = ? and id > ? and is_del = 0", userId, lastId).
		Limit(pageSize).Find(&res).Error
	return res, err
}

func (m *SlSlUserShortURLDao) UpdateByShortURL(shortURL string, data map[string]any, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlUserShortURL{}).TableName()).Where("short_url = ?", shortURL).Updates(data).Error
}

func (m *SlSlUserShortURLDao) DeleteByShortURL(shortURL string, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlUserShortURL{}).TableName()).Where("short_url = ?", shortURL).Delete(&SlUserShortURL{}).Error
}
