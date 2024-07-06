package db

import (
	"context"
	"errors"
	"short-link/database/mysql"
	"time"

	"gorm.io/gorm"
)

// 黑名单
type SlBlackList struct {
	Id        uint64 `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	ShortURL  string `gorm:"column:short_url;type:varchar(50);comment:短链;NOT NULL" json:"shortUrl"`
	IP        uint32 `gorm:"column:IP;type:bigint(20) unsigned;comment:IP" json:"IP"`
	CreatedAt int64  `gorm:"column:created_at;type:bigint(20);NOT NULL" json:"createdAt"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint(20);NOT NULL" json:"updatedAt"`
}

func (m *SlBlackList) TableName() string {
	return "sl_black_list"
}

type SlBlackListDao struct {
	db *gorm.DB
}

func NewSlBlackListDao(ctx context.Context, db ...*gorm.DB) *SlBlackListDao {
	client := mysql.NewDBClient(ctx)
	if len(db) > 0 {
		client = db[0]
	}
	return &SlBlackListDao{db: client}
}

func (m *SlBlackList) BeforeCreate(_ *gorm.DB) (err error) {
	m.CreatedAt = time.Now().UnixMilli()
	m.UpdatedAt = time.Now().UnixMilli()
	return
}

func (m *SlBlackListDao) Create(l *SlBlackList, db ...*gorm.DB) error {
	tx := m.db
	if len(db) > 0 {
		tx = db[0]
	}
	return tx.Table((&SlBlackList{}).TableName()).Create(l).Error
}

func (m *SlBlackListDao) GetByShortURL(shortURL string) ([]*SlBlackList, error) {
	var res []*SlBlackList
	err := m.db.Table((&SlBlackList{}).TableName()).
		Where("short_url = ?", shortURL).
		Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *SlBlackListDao) List(shortURL string, ip uint32, page, pageSize int) ([]*SlBlackList, int64, error) {
	var res []*SlBlackList
	var count int64
	query := m.db.Table((&SlBlackList{}).TableName()).
		Where("short_url = ?", shortURL)
	if ip > 0 {
		query = query.Where("IP = ?", ip)
	}
	err := query.Order("id DESC").
		Count(&count).
		Offset((page - 1) * pageSize).
		Limit(pageSize).Find(&res).Error
	return res, count, err
}

func (m *SlBlackListDao) Delete(id uint64) error {
	err := m.db.Table((&SlBlackList{}).TableName()).
		Where("id = ?", id).
		Delete(&SlBlackList{}).Error

	return err
}

func (m *SlBlackListDao) DeleteByshortURL(shortURL string) error {
	err := m.db.Table((&SlBlackList{}).TableName()).
		Where("short_url = ?", shortURL).
		Delete(&SlBlackList{}).Error
	return err
}

func (m *SlBlackListDao) GetByID(id uint64) (*SlBlackList, error) {
	var res SlBlackList
	err := m.db.Table((&SlBlackList{}).TableName()).
		Where("id = ?", id).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, err
}
