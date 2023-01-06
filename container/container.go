package container

import (
	"flag"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/config"
	_ "github.com/br0-space/bot/pkg/config"
	"github.com/br0-space/bot/pkg/db"
	"github.com/br0-space/bot/pkg/logger"
	"github.com/br0-space/bot/pkg/matcher"
	"github.com/br0-space/bot/pkg/repo"
	"github.com/br0-space/bot/pkg/telegram"
	"github.com/br0-space/bot/pkg/webhook"
	"gorm.io/gorm"
)

var loggerInstance interfaces.LoggerInterface
var configInstance *interfaces.ConfigStruct

func runsAsTest() bool {
	return flag.Lookup("test.v") != nil
}

func ProvideLogger() interfaces.LoggerInterface {
	if loggerInstance == nil {
		if runsAsTest() {
			loggerInstance = logger.NewNullLogger()
		} else {
			loggerInstance = logger.NewDefaultLogger()
		}
	}
	return loggerInstance
}

func ProvideConfig() *interfaces.ConfigStruct {
	if configInstance == nil {
		if runsAsTest() {
			configInstance = config.NewTestConfig()
		} else {
			configInstance = config.NewConfig()
		}
	}
	return configInstance
}

func ProvideMatchersRegistry() interfaces.MatcherRegistryInterface {
	return matcher.NewRegistry(
		ProvideLogger(),
		ProvideTelegramClient(),
		ProvideDatabaseRepository(),
	)
}

func ProvideTelegramWebhookHandler() interfaces.TelegramWebhookHandlerInterface {
	return webhook.NewHandler(
		ProvideLogger(),
		ProvideConfig(),
		ProvideMatchersRegistry(),
	)
}

func ProvideTelegramWebhookTools() interfaces.TelegramWebhookToolsInterface {
	if runsAsTest() {
		return webhook.NewMockTools()
	} else {
		return webhook.NewProdTools(
			ProvideLogger(),
			ProvideConfig().Telegram,
		)
	}
}

func ProvideTelegramClient() interfaces.TelegramClientInterface {
	if runsAsTest() {
		return telegram.NewMockClient()
	} else {
		return telegram.NewProdClient(
			ProvideLogger(),
			ProvideConfig().Telegram,
		)
	}
}

func ProvideDatabaseConnection() *gorm.DB {
	return db.NewConnection(
		ProvideLogger(),
		ProvideConfig().Database,
	)
}

func ProvideDatabaseMigration() *db.DatabaseMigration {
	return db.NewDatabaseMigration(
		ProvideLogger(),
		ProvideDatabaseRepository(),
	)
}

func ProvideDatabaseRepository() interfaces.DatabaseRepositoryInterface {
	tx := ProvideDatabaseConnection()

	return repo.NewRepository(
		ProvideLogger(),
		ProvideMessageStatsRepo(tx),
		ProvidePlusplusRepo(tx),
		ProvideStatsRepo(tx),
	)
}

func ProvideMessageStatsRepo(tx *gorm.DB) interfaces.MessageStatsRepoInterface {
	return repo.NewMessageStatsRepo(
		ProvideLogger(),
		tx,
	)
}

func ProvidePlusplusRepo(tx *gorm.DB) interfaces.PlusplusRepoInterface {
	return repo.NewPlusplusRepo(
		ProvideLogger(),
		tx,
	)
}

func ProvideStatsRepo(tx *gorm.DB) interfaces.StatsRepoInterface {
	return repo.NewStatsRepo(
		ProvideLogger(),
		tx,
	)
}
