package container

import (
	"flag"

	"github.com/br0-space/bot/internal/config"
	"github.com/br0-space/bot/internal/logger"
	"github.com/br0-space/bot/internal/telegram"
	"github.com/br0-space/bot/internal/telegram/client"
)

func ProvideConfig() *config.Config {
	if runsAsTest() {
		return config.NewTestConfig()
	}
	return config.NewLiveConfig()
}

func ProvideLoggerService() logger.Interface {
	if runsAsTest() {
		return logger.NewNullLogger()
	}
	return logger.NewLiveLogger()
}

func ProvideTelegramService() telegram.Interface {
	// if runsAsTest() {
		return telegram.NewMockTelegram()
	// }
	// return nil
}

func ProvideTelegramClientService() *client.Client {
	return client.NewClient()
}

func runsAsTest() bool {
	return flag.Lookup("test.v") != nil
}