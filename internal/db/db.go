package db

import (
	"fmt"
	"log"

	"github.com/neovg/kmptnzbot/internal/config"
	"github.com/neovg/kmptnzbot/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	logger.Log.Debug("init database")

	var db *gorm.DB
	var err error
	switch config.Cfg.Database.Driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Cfg.Database.SQLite.File), &gorm.Config{})
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s TimeZone=%s",
			config.Cfg.Database.Postgres.Host,
			config.Cfg.Database.Postgres.Port,
			config.Cfg.Database.Postgres.DBName,
			config.Cfg.Database.Postgres.User,
			config.Cfg.Database.Postgres.Password,
			config.Cfg.Database.Postgres.SSL,
			config.Cfg.Database.Postgres.Timezone,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		log.Panicln("unknown database driver", config.Cfg.Database.Driver)
	}

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	MigratePlusplus(db)
	MigrateStats(db)
	MigrateMessageStats(db)

	DB = db
}
