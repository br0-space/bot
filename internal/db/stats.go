package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/br0-space/bot/internal/oldtelegram"
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

// Is called when any new message arrives update updates that users stats
func UpdateStats(user oldtelegram.User) {
	// Prevent race conditions
	mutexStats.Lock()

	// Allow other goroutines to execute this block again if the function is finished
	defer mutexStats.Unlock()

	// Try to insert a new entry
	// If the insert fails, update the existing entry instead (upsert)
	DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"username":  user.UsernameOrName(),
			"posts":     gorm.Expr("stats.posts + 1"),
			"last_post": time.Now(),
		}),
	}).Create(&Stats{
		UserID:   user.ID,
		Username: user.UsernameOrName(),
		Posts:    1,
		LastPost: time.Now(),
	})
}

// Finds all stats sorted by the number of posts
func FindStatsTop() []Stats {
	var records []Stats
	DB.Order("posts desc").Find(&records)

	return records
}

// Finds all usernames for use in @all
func FindAllUsernames(exclude string) []string {
	var records []Stats
	DB.Where(fmt.Sprintf("username like '@%%' and username != '@%s'", exclude)).Order("username asc").Find(&records)

	usernames := make([]string, 0, len(records))
	for _, record := range records {
		usernames = append(usernames, record.Username)
	}

	return usernames
}
