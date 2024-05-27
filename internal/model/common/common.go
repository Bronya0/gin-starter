package common

import (
	"time"
)

type BaseModel struct {
	ID          uint      `gorm:"primarykey"` // 主键ID
	CreateTime  time.Time // 创建时间
	UpdatedTime time.Time // 更新时间
}
