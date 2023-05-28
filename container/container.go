package container

import (
	"flag"
	"sync"

	logger "github.com/br0-space/bot-logger"
	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/config"
	"github.com/br0-space/bot/pkg/db"
	"github.com/br0-space/bot/pkg/fortune"
	"github.com/br0-space/bot/pkg/matchers/atall"
	"github.com/br0-space/bot/pkg/matchers/buzzwords"
	"github.com/br0-space/bot/pkg/matchers/choose"
	fortune2 "github.com/br0-space/bot/pkg/matchers/fortune"
	"github.com/br0-space/bot/pkg/matchers/goodmorning"
	"github.com/br0-space/bot/pkg/matchers/janein"
	"github.com/br0-space/bot/pkg/matchers/ping"
	"github.com/br0-space/bot/pkg/matchers/plusplus"
	"github.com/br0-space/bot/pkg/matchers/stats"
	"github.com/br0-space/bot/pkg/matchers/topflop"
	xkcd2 "github.com/br0-space/bot/pkg/matchers/xkcd"
	"github.com/br0-space/bot/pkg/repo"
	"github.com/br0-space/bot/pkg/state"
	"github.com/br0-space/bot/pkg/xkcd"
	"gorm.io/gorm"
)

var (
	configInstance          *interfaces.ConfigStruct
	configLock              = &sync.Mutex{}
	matcherRegistryInstance *matcher.Registry
	matcherRegistryLock     = &sync.Mutex{}
	stateInstance           interfaces.StateServiceInterface
	stateLock               = &sync.Mutex{}
)

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

func ProvideMatchersRegistry() *matcher.Registry {
	matcherRegistryLock.Lock()
	defer matcherRegistryLock.Unlock()

	if matcherRegistryInstance == nil {
		matcherRegistryInstance = matcher.NewRegistry(
			ProvideLogger(),
			ProvideTelegramClient(),
		)
		matcherRegistryInstance.Register(atall.MakeMatcher(ProvideUserStatsRepo()))
		matcherRegistryInstance.Register(buzzwords.MakeMatcher(ProvidePlusplusRepo()))
		matcherRegistryInstance.Register(choose.MakeMatcher())
		matcherRegistryInstance.Register(goodmorning.MakeMatcher(ProvideState(), ProvideFortuneService()))
		matcherRegistryInstance.Register(fortune2.MakeMatcher(ProvideFortuneService()))
		matcherRegistryInstance.Register(janein.MakeMatcher())
		matcherRegistryInstance.Register(ping.MakeMatcher())
		matcherRegistryInstance.Register(plusplus.MakeMatcher(ProvidePlusplusRepo()))
		matcherRegistryInstance.Register(stats.MakeMatcher(ProvideUserStatsRepo()))
		matcherRegistryInstance.Register(topflop.MakeMatcher(ProvidePlusplusRepo()))
		matcherRegistryInstance.Register(xkcd2.MakeMatcher(ProvideXkcdService()))
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

func ProvideXkcdService() interfaces.XkcdServiceInterface {
	return xkcd.MakeService()
}
