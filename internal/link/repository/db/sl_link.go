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
	ID        int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	OriginUrl string `gorm:"column:origin_url;NOT NULL"`           // 原始链接
	ShortUrl  string `gorm:"column:short_url;NOT NULL"`            // 短链
	ExpiredAt int64  `gorm:"column:expired_at;default:0;NOT NULL"` // 过期时间
	UserId    uint64 `gorm:"column:user_id;default:0;NOT NULL"`    // 用户ID
	CreatedAt int64  `gorm:"column:created_at;default:0;NOT NULL"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;default:0;NOT NULL"` // 更新时间
	IsDel     int    `gorm:"column:is_del;default:0;NOT NULL"`
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

func (m *SlLinkDao) Create(l *SlLink) error {
	return m.db.Table((&SlLink{}).TableName(l.ShortUrl)).Create(l).Error
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
