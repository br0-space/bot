package db

import (
	"fmt"
	"time"

	logger "github.com/br0-space/bot-logger"
	"github.com/br0-space/bot/interfaces"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewConnection(logger logger.Interface, config interfaces.DatabaseConfigStruct) *gorm.DB {
	var dialector gorm.Dialector

	switch config.Driver {
	case "sqlite":
		dialector = sqlite.Open(config.SQLite.File)
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s TimeZone=%s",
			config.PostgreSQL.Host,
			config.PostgreSQL.Port,
			config.PostgreSQL.DBName,
			config.PostgreSQL.User,
			config.PostgreSQL.Password,
			config.PostgreSQL.SSL,
			config.PostgreSQL.Timezone,
		)
		dialector = postgres.Open(dsn)
	default:
		logger.Panic("unknown database driver", config.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		FullSaveAssociations:   false,
		AllowGlobalUpdate:      false,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 NewGormLoggerBridge(logger),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		logger.Error("failed to connect database:", err)

		return nil
	}

	return db
}
