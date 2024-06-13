package db

type SlUserShortUrl struct {
	ID        int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ShortUrl  string `gorm:"column:short_url;NOT NULL"`            // 短链
	UserId    uint64 `gorm:"column:user_id;default:0;NOT NULL"`    // 用户ID
	CreatedAt int64  `gorm:"column:created_at;default:0;NOT NULL"` // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;default:0;NOT NULL"` // 更新时间
	IsDel     int    `gorm:"column:is_del;default:0;NOT NULL"`
}

func (m *SlUserShortUrl) TableName() string {
	return "sl_user_short_url"
}
