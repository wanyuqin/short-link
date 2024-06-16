package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"short-link/database/mysql"
	"time"
)

type SlUser struct {
	ID        uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Username  string `gorm:"column:username;NOT NULL"`             // 用户名
	Password  string `gorm:"column:password;NOT NULL"`             // 密码
	Salt      string `gorm:"column:salt;NOT NULL"`                 // 加密盐
	CreatedAt int64  `gorm:"column:created_at;default:0;NOT NULL"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;default:0;NOT NULL"` // 更新时间
	IsDel     int    `gorm:"column:is_del;default:0;NOT NULL"`
}

func (m *SlUser) TableName() string {
	return "sl_user"
}

type SlUserDao struct {
	db *gorm.DB
}

func NewSlUserDao(ctx context.Context, db ...*gorm.DB) *SlUserDao {
	client := mysql.NewDBClient(ctx)
	if len(db) > 0 {
		client = db[0]
	}
	return &SlUserDao{
		db: client,
	}
}

func (m *SlUser) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().UnixMilli()
	m.UpdatedAt = time.Now().UnixMilli()
	return
}

func (m *SlUserDao) Create(u *SlUser) error {
	return m.db.Create(&u).Error
}

func (m *SlUserDao) GetByUname(uname string) (*SlUser, error) {
	var u SlUser
	err := m.db.Table((&SlUser{}).TableName()).Where("username = ?", uname).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return &u, err
}

func (m *SlUserDao) GetByUnameAndPwd(uname string, pwd string) (*SlUser, error) {
	var u SlUser
	err := m.db.Table((&SlUser{}).TableName()).Where("username = ? and password = ?", uname, pwd).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return &u, err
}
