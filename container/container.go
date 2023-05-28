package container

import (
	"flag"
	logger "github.com/br0-space/bot-logger"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/config"
	"github.com/br0-space/bot/pkg/db"
	"github.com/br0-space/bot/pkg/fortune"
	"github.com/br0-space/bot/pkg/matcher"
	"github.com/br0-space/bot/pkg/repo"
	"github.com/br0-space/bot/pkg/songlink"
	"github.com/br0-space/bot/pkg/state"
	"github.com/br0-space/bot/pkg/xkcd"
	"gorm.io/gorm"
	"sync"
)

var configInstance *interfaces.ConfigStruct
var configLock = &sync.Mutex{}
var matcherRegistryInstance interfaces.MatcherRegistryInterface
var matcherRegistryLock = &sync.Mutex{}
var stateInstance interfaces.StateServiceInterface
var stateLock = &sync.Mutex{}

func runsAsTest() bool {
	return flag.Lookup("test.v") != nil
}

func ProvideLogger() logger.Interface {
	return logger.New()
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
	matcherRegistryLock.Lock()
	defer matcherRegistryLock.Unlock()

	if matcherRegistryInstance == nil {
		matcherRegistryInstance = matcher.NewRegistry(
			ProvideState(),
			ProvideTelegramClient(),
			ProvideMessageStatsRepo(),
			ProvidePlusplusRepo(),
			ProvideUserStatsRepo(),
			ProvideFortuneService(),
			ProvideSonglinkService(),
			ProvideXkcdService(),
		)
	}

	return matcherRegistryInstance
}

func ProvideState() interfaces.StateServiceInterface {
	stateLock.Lock()
	defer stateLock.Unlock()

	if stateInstance == nil {
		stateInstance = state.NewService(
			ProvideUserStatsRepo(),
			ProvideMessageStatsRepo(),
		)
	}

	return stateInstance
}

func ProvideTelegramWebhookHandler() telegramclient.WebhookHandlerInterface {
	matchersRegistry := ProvideMatchersRegistry()
	stateService := ProvideState()

	return telegramclient.NewHandler(
		&ProvideConfig().Telegram,
		func(messageIn telegramclient.WebhookMessageStruct) {
			matchersRegistry.Process(messageIn)
			stateService.ProcessMessage(messageIn)
		},
	)
}

func ProvideTelegramClient() telegramclient.ClientInterface {
	if runsAsTest() {
		return telegramclient.NewMockClient()
	} else {
		return telegramclient.NewClient(
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
		ProvideMessageStatsRepo(),
		ProvidePlusplusRepo(),
		ProvideUserStatsRepo(),
	)
}

func ProvideMessageStatsRepo() interfaces.MessageStatsRepoInterface {
	return repo.NewMessageStatsRepo(
		ProvideDatabaseConnection(),
	)
}

func ProvidePlusplusRepo() interfaces.PlusplusRepoInterface {
	return repo.NewPlusplusRepo(
		ProvideDatabaseConnection(),
	)
}

func ProvideUserStatsRepo() interfaces.UserStatsRepoInterface {
	return repo.NewUserStatsRepo(
		ProvideDatabaseConnection(),
	)
}

func ProvideFortuneService() fortune.Service {
	return fortune.MakeService()
}

func ProvideSonglinkService() interfaces.SonglinkServiceInterface {
	return songlink.MakeService()
}

func ProvideXkcdService() interfaces.XkcdServiceInterface {
	return xkcd.MakeService()
}
