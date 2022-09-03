package db

import (
	"gorm.io/gorm"
	"time"
)

type Stats struct {
	gorm.Model
	ChatID   int64  `gorm:"<-:create;index:idx_stats_chat_id_user_id,unique"`
	UserID   int64  `gorm:"<-:create;index:idx_stats_chat_id_user_id,unique"`
	Username string `gorm:"<-"`
	Posts    uint32 `gorm:"<-"`
	LastPost time.Time
}
