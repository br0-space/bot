package db

import (
	"time"

	"github.com/br0-space/bot/internal/telegram"
	"gorm.io/gorm"
)

type MessageStats struct {
	gorm.Model
	UserID    int64     `gorm:"index;<-:create"`
	UserStats Stats     `gorm:"foreignKey:user_id;references:user_id;constraint:OnDelete:CASCADE"`
	Time      time.Time `gorm:"index;<-:create"`
	Words     int       `gorm:"<-:create"`
}

type WordCount struct {
	UserID   int64
	Username string
	Words    int
}

// Migrates the table message_stats
func MigrateMessageStats(db *gorm.DB) {
	if err := db.AutoMigrate(&MessageStats{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

func InsertMessageStats(requestMessage telegram.RequestMessage) {
	DB.Create(&MessageStats{
		UserID: requestMessage.From.ID,
		Time:   time.Now(),
		Words:  requestMessage.WordCount(),
	})
}

func GetWordCounts() []WordCount {
	var records []WordCount
	DB.Model(&MessageStats{}).Joins("UserStats").Select(`"message_stats".user_id, "UserStats".username, count("message_stats".words) as words`).Group(`"message_stats".user_id, "UserStats".id, "UserStats".username`).Order(`count("message_stats".words) desc`).Scan(&records)

	return records
}
