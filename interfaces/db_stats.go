package interfaces

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

type StatsUserStruct struct {
	ID       int64
	Username string
	Posts    uint32
	LastPost time.Time
}

type StatsRepoInterface interface {
	UpdateStats(chatID int64, userID int64, username string) error
	GetKnownUsers(chatID int64) ([]StatsUserStruct, error)
	GetTopUsers(chatID int64) ([]StatsUserStruct, error)
}
