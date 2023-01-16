package container

import (
	"flag"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/config"
	_ "github.com/br0-space/bot/pkg/config"
	"github.com/br0-space/bot/pkg/db"
	"github.com/br0-space/bot/pkg/fortune"
	"github.com/br0-space/bot/pkg/logger"
	"github.com/br0-space/bot/pkg/matcher"
	"github.com/br0-space/bot/pkg/repo"
	"github.com/br0-space/bot/pkg/songlink"
	"github.com/br0-space/bot/pkg/state"
	"github.com/br0-space/bot/pkg/telegram"
	"github.com/br0-space/bot/pkg/webhook"
	"gorm.io/gorm"
	"sync"
)

var loggerInstance interfaces.LoggerInterface
var loggerLock = &sync.Mutex{}
var configInstance *interfaces.ConfigStruct
var configLock = &sync.Mutex{}
var stateInstance interfaces.StateServiceInterface
var stateLock = &sync.Mutex{}

func runsAsTest() bool {
	return flag.Lookup("test.v") != nil
}

func ProvideLogger() interfaces.LoggerInterface {
	loggerLock.Lock()
	defer loggerLock.Unlock()

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
	configLock.Lock()
	defer configLock.Unlock()

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
		ProvideState(),
		ProvideTelegramClient(),
		ProvideMessageStatsRepo(),
		ProvidePlusplusRepo(),
		ProvideUserStatsRepo(),
		ProvideFortuneService(),
		ProvideSonglinkService(),
	)
}

func ProvideState() interfaces.StateServiceInterface {
	stateLock.Lock()
	defer stateLock.Unlock()

	if stateInstance == nil {
		stateInstance = state.NewService(
			ProvideLogger(),
			ProvideUserStatsRepo(),
			ProvideMessageStatsRepo(),
		)
	}

	return stateInstance
}

func ProvideTelegramWebhookHandler() interfaces.TelegramWebhookHandlerInterface {
	return webhook.NewHandler(
		ProvideLogger(),
		ProvideConfig(),
		ProvideMatchersRegistry(),
		ProvideState(),
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

func ProvideDatabaseMigration() interfaces.DatabaseMigrationInterface {
	return db.MakeDatabaseMigration(
		ProvideLogger(),
		ProvideMessageStatsRepo(),
		ProvidePlusplusRepo(),
		ProvideUserStatsRepo(),
	)
}

func ProvideMessageStatsRepo() interfaces.MessageStatsRepoInterface {
	return repo.NewMessageStatsRepo(
		ProvideLogger(),
		ProvideDatabaseConnection(),
	)
}

func ProvidePlusplusRepo() interfaces.PlusplusRepoInterface {
	return repo.NewPlusplusRepo(
		ProvideLogger(),
		ProvideDatabaseConnection(),
	)
}

func ProvideUserStatsRepo() interfaces.UserStatsRepoInterface {
	return repo.NewUserStatsRepo(
		ProvideLogger(),
		ProvideDatabaseConnection(),
	)
}

func ProvideFortuneService() fortune.Service {
	return fortune.MakeService()
}

func ProvideSonglinkService() interfaces.SonglinkServiceInterface {
	return songlink.MakeService()
}
