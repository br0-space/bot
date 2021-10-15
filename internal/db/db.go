package db

import (
	"fmt"
	"log"

	"github.com/br0-space/bot/internal/oldconfig"
	"github.com/br0-space/bot/internal/oldlogger"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	oldlogger.Log.Debug("init database")

	var db *gorm.DB
	var err error
	switch oldconfig.Cfg.Database.Driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(oldconfig.Cfg.Database.SQLite.File), &gorm.Config{})
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s TimeZone=%s",
			oldconfig.Cfg.Database.Postgres.Host,
			oldconfig.Cfg.Database.Postgres.Port,
			oldconfig.Cfg.Database.Postgres.DBName,
			oldconfig.Cfg.Database.Postgres.User,
			oldconfig.Cfg.Database.Postgres.Password,
			oldconfig.Cfg.Database.Postgres.SSL,
			oldconfig.Cfg.Database.Postgres.Timezone,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		log.Panicln("unknown database driver", oldconfig.Cfg.Database.Driver)
	}

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	MigratePlusplus(db)
	MigrateStats(db)
	MigrateMessageStats(db)

	DB = db
}
