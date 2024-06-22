package db

import (
	"context"
	"crypto/md5"
	"errors"
	"gorm.io/gorm"
	"short-link/database/mysql"
	"strconv"
	"time"
)

type SlLink struct {
	ID        uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"ID,omitempty"`
	OriginUrl string `gorm:"column:origin_url;NOT NULL" json:"originUrl,omitempty"`           // 原始链接
	ShortUrl  string `gorm:"column:short_url;NOT NULL" json:"shortUrl,omitempty"`             // 短链
	ExpiredAt int64  `gorm:"column:expired_at;default:0;NOT NULL" json:"expiredAt,omitempty"` // 过期时间
	UserId    uint64 `gorm:"column:user_id;default:0;NOT NULL" json:"userId,omitempty"`       // 用户ID
	CreatedAt int64  `gorm:"column:created_at;default:0;NOT NULL" json:"createdAt,omitempty"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;default:0;NOT NULL" json:"updatedAt,omitempty"` // 更新时间
	IsDel     int    `gorm:"column:is_del;default:0;NOT NULL" json:"isDel,omitempty"`
}

func (m *SlLink) TableName(shortUrl string) string {
	// 计算 shortUrl 的 MD5 哈希值
	hash := md5.Sum([]byte(shortUrl))

	// 取哈希值的最后一位，将其转换为 0-9 之间的整数
	index := hash[len(hash)-1] % 10

	// 构造表名
	tableName := "sl_link_" + strconv.Itoa(int(index))

	return tableName
}

type SlLinkDao struct {
	db *gorm.DB
}

func NewSlLinkDao(ctx context.Context, db ...*gorm.DB) *SlLinkDao {
	client := mysql.NewDBClient(ctx)
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

func (m *SlLinkDao) Create(l *SlLink, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlLink{}).TableName(l.ShortUrl)).Create(l).Error
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

func (m *SlLinkDao) GetByShortUrlsWithTableName(tableName string, shortUrls []string) ([]*SlLink, error) {
	var res []*SlLink
	err := m.db.Table(tableName).Where("short_url in (?) and is_del = 0", shortUrls).
		Order("created_at DESC").
		Find(&res).Error
	return res, err
}

func (m *SlLinkDao) UpdateByShortUrl(shortUrl string, data map[string]interface{}, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlLink{}).TableName(shortUrl)).Where("short_url = ?", shortUrl).Updates(data).Error
}

func (m *SlLinkDao) DeleteByShort(shortUrl string, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlLink{}).TableName(shortUrl)).Where("short_url = ?", shortUrl).Delete(&SlLink{}).Error
}

func (m *SlLinkDao) GetByShortUrl(shortUrl string) (*SlLink, error) {
	var res SlLink
	err := m.db.Table((&SlLink{}).TableName(shortUrl)).Where("short_url = ? and is_del = ?", shortUrl, 0).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}
