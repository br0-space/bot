package interfaces

import (
	"time"

	"gorm.io/gorm"
)

type MessageStats struct {
	gorm.Model `exhaustruct:"optional"`
	UserID     int64 `gorm:"<-:create;index"`
	// UserStats Stats     `gorm:"foreignKey:user_id;references:user_id;constraint:OnDelete:CASCADE"`
	Time  time.Time `gorm:"<-:create;index"`
	Words int       `gorm:"<-:create"`
}

type MessageStatsWordCountStruct struct {
	UserID   int64
	Username string
	Words    int
}

type MessageStatsRepoInterface interface {
	InsertMessageStats(userID int64, words int) error
	GetKnownUserIDs() ([]int64, error)
	GetWordCounts() ([]MessageStatsWordCountStruct, error)
}
