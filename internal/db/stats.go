package db

import (
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Model for table stats
type Stats struct {
	gorm.Model
	UserID   int64  `gorm:"uniqueIndex;<-:create"`
	Username string `gorm:"<-"`
	Posts    uint32 `gorm:"<-"`
	LastPost time.Time
}

var mutexStats sync.Mutex

// Migrates the table stats
func MigrateStats(db *gorm.DB) {
	if err := db.AutoMigrate(&Stats{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

func UpdateStats(userID int64, username string) {
	// Prevent race conditions
	mutexStats.Lock()

	// Allow other goroutines to execute this block again if the function is finished
	defer mutexStats.Unlock()

	// Try to insert a new entry
	// If the insert fails, update the existing entry instead (upsert)
	DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"username": username,
			"posts": gorm.Expr("posts + 1"),
			"last_post": time.Now(),
		}),
	}).Create(&Stats{
		UserID:   userID,
		Username: username,
		Posts:    1,
		LastPost: time.Now(),
	})
}

func GetStatsTops() []Stats {
	var records []Stats
	DB.Order("last_post desc").Find(&records)

	return records
}