package telegram

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/neovg/kmptnzbot/internal/config"
	"github.com/neovg/kmptnzbot/internal/logger"
)

// Use the Telegram API to set the webhook URL messages to the bot will be sent to
func SetWebhookURL() {
	// Only set the webhook URL when it is set in .env or an environment variable
	if len(config.Cfg.Telegram.WebhookURL) > 0 {
		logger.Log.Info("set webhook url to", config.Cfg.Telegram.WebhookURL)

		apiUrl := fmt.Sprintf(config.Cfg.Telegram.BaseUrl, config.Cfg.Telegram.ApiKey) + config.Cfg.Telegram.EndpointSetWebhook
		_, err := http.PostForm(apiUrl, url.Values{
			"url": {config.Cfg.Telegram.WebhookURL},
		})
		if err != nil {
			logger.Log.Panic("cannot set Telegram webhook URL:", err)
		}
	}
}