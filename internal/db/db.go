package db

import (
	"github.com/neovg/kmptnzbot/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	logger.Log.Debug("init database")
    dsn := "host=db user=gorm password=example dbname=exmaple port=5132 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	MigratePlusplus(db)

	DB = db
}
