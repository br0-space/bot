package oldtelegram

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/br0-space/bot/internal/oldconfig"
	"github.com/br0-space/bot/internal/oldlogger"
)

// Use the telegram API to set the webhook URL messages to the bot will be sent to
func SetWebhookURL() {
	// Only set the webhook URL when it is set in .env or an environment variable
	if len(oldconfig.Cfg.Telegram.WebhookURL) > 0 {
		oldlogger.Log.Info("set webhook url to", oldconfig.Cfg.Telegram.WebhookURL)

		apiUrl := fmt.Sprintf(oldconfig.Cfg.Telegram.BaseUrl, oldconfig.Cfg.Telegram.ApiKey) + oldconfig.Cfg.Telegram.EndpointSetWebhook
		_, err := http.PostForm(apiUrl, url.Values{
			"url": {oldconfig.Cfg.Telegram.WebhookURL},
		})
		if err != nil {
			oldlogger.Log.Panic("cannot set telegram webhook URL:", err)
		}
	}
}
