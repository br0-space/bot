package interfaces

import (
	"gorm.io/gorm"
	"time"
)

type Stats struct {
	gorm.Model
	UserID   int64  `gorm:"<-:create;uniqueIndex"`
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
	UpdateStats(userID int64, username string) error
	GetKnownUsers() ([]StatsUserStruct, error)
	GetTopUsers() ([]StatsUserStruct, error)
}
