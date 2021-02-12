package db

import (
	"github.com/neovg/kmptnzbot/internal/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	logger.Log.Debug("init database")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	MigratePlusplus(db)
	MigrateStats(db)

	DB = db
}