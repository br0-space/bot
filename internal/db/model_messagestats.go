package db

import (
	"gorm.io/gorm"
	"time"
)

type MessageStats struct {
	gorm.Model
	ChatID int64 `gorm:"<-:create;index"`
	UserID int64 `gorm:"<-:create;index"`
	//UserStats Stats     `gorm:"foreignKey:user_id;references:user_id;constraint:OnDelete:CASCADE"`
	Time  time.Time `gorm:"<-:create;index"`
	Words int       `gorm:"<-:create"`
}
