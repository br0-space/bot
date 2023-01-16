package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"net/http"
)

type ProdClient struct {
	Log interfaces.LoggerInterface
	Cfg interfaces.TelegramConfigStruct
}

type sendMessageResponse struct {
	Ok          bool   `json:"ok"`
	Result      bool   `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func NewProdClient(logger interfaces.LoggerInterface, config interfaces.TelegramConfigStruct) *ProdClient {
	return &ProdClient{
		Log: logger,
		Cfg: config,
	}
}

func (c ProdClient) SendMessage(chatID int64, message interfaces.TelegramMessageStruct) error {
	switch {
	case message.Photo != "":
		c.Log.Debugf("Sending photo: %s", message.Photo)
	default:
		c.Log.Debugf("Sending message: %s", message.Text)
	}

	message.ChatID = chatID

	url := c.getUrl(message)
	requestBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	c.Log.Debugf("Sending POST request to %s", url)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusOK {
		c.Log.Debug("Successfully sent message to Telegram")

		return nil
	}

	responseBody := &sendMessageResponse{}
	if err = json.NewDecoder(response.Body).Decode(responseBody); err != nil {
		return fmt.Errorf("SendMessage failed with %s: unable to decode response body", response.Status)
	}
	return fmt.Errorf("SendMessage failed with %d: %s", responseBody.ErrorCode, responseBody.Description)
}

func (c ProdClient) getUrl(message interfaces.TelegramMessageStruct) string {
	switch {
	case message.Photo != "":
		return fmt.Sprintf(c.Cfg.BaseUrl, c.Cfg.ApiKey) + c.Cfg.EndpointSendPhoto
	default:
		return fmt.Sprintf(c.Cfg.BaseUrl, c.Cfg.ApiKey) + c.Cfg.EndpointSendMessage
	}
}
