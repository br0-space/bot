package webhook

import (
	"encoding/json"
	"fmt"
	logger "github.com/br0-space/bot-logger"
	"github.com/br0-space/bot/interfaces"
	"net/http"
	"net/url"
)

type ProdTools struct {
	Log logger.Interface
	Cfg interfaces.TelegramConfigStruct
}

type setWebhookURLResponse struct {
	Ok          bool   `json:"ok"`
	Result      bool   `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func NewProdTools(config interfaces.TelegramConfigStruct) *ProdTools {
	return &ProdTools{
		Log: logger.New(),
		Cfg: config,
	}
}

func (t *ProdTools) SetWebhookURL() {
	if t.Cfg.WebhookURL == "" {
		t.Log.Info("Not setting Telegram webhook URL")
		return
	}

	t.Log.Info("Setting Telegram webhook URL to", t.Cfg.WebhookURL)

	apiUrl := fmt.Sprintf(t.Cfg.BaseUrl, t.Cfg.ApiKey) + t.Cfg.EndpointSetWebhook

	t.Log.Debug("Sending POST request to", apiUrl)

	if resp, err := http.PostForm(apiUrl, url.Values{
		"url": {t.Cfg.WebhookURL},
	}); err != nil {
		t.Log.Panic("Unable to set Telegram webhook URL:", err)
	} else {
		body := &setWebhookURLResponse{}
		if err = json.NewDecoder(resp.Body).Decode(body); err != nil {
			t.Log.Fatal("Unable to decode response body:", err)
		}

		if !body.Ok {
			t.Log.Fatal("Unable to set Telegram webhook URL:", body.Description)
		}

		t.Log.Debug("Successfully set Telegram webhook URL")
	}
}
