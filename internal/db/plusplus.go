package db

import (
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Model for table plusplus
type Plusplus struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex;<-:create"`
	Value int    `gorm:"index;<-"`
}

var mutexPlusplus sync.Mutex

// Migrates the table plusplus
func MigratePlusplus(db *gorm.DB) {
	if err := db.AutoMigrate(&Plusplus{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

// Atomically increments a plusplus entry and returns the new value
func IncrementPlusplus(name string, increment int) int {
	// Prevent race conditions between two plusplus matcher goroutines trying to increment the same entry at the same time
	mutexPlusplus.Lock()

	// Allow other goroutines to execute this block again if the function is finished
	defer mutexPlusplus.Unlock()

	// Try to insert a new entry
	// If the insert fails, update the existing entry instead (upsert)
	DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"value": gorm.Expr("plusplus.value + ?", increment),
		}),
	}).Create(&Plusplus{
		Name:  name,
		Value: increment,
	})

	// Read the new value (guaranteed to be the result of the insert by the mutex)
	var record Plusplus
	DB.Where("name = ?", name).First(&record)

	return record.Value
}

// Find the top 10 plusplus entries
func FindPlusplusTops() []Plusplus {
	var records []Plusplus
	DB.Where("value > 0").Order("value desc").Limit(10).Find(&records)

	return records
}

// Find the lowest 10 plusplus entries
func FindPlusplusFlops() []Plusplus {
	var records []Plusplus
	DB.Where("value <= 0").Order("value asc").Limit(10).Find(&records)

	return records
}
