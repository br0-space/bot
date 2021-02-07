package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/neovg/kmptnzbot/internal/config"
	"github.com/neovg/kmptnzbot/internal/logger"
)

// Create a struct that is accepted by Telegram's sendMessage endpoint
// https://core.telegram.org/bots/api#sendmessage

type Message struct {
	ChatID                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotification   bool   `json:"disable_notification"`
	ReplyToMessageID      int64  `json:"reply_to_message_id"`
}

type responseBody struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int16  `json:"error_code"`
	Description string `json:"description"`
}

// Send a message to the telegram channel the request message came from
func SendMessage(requestMessage RequestMessage, responseMessage Message) error {
	logger.Log.Infof("send to %d: %s", requestMessage.Chat.ID, responseMessage.Text)

	// Send the message to the same chat where the request message came from
	responseMessage.ChatID = requestMessage.Chat.ID

	// Request body will be JSON
	requestBytes, err := json.Marshal(responseMessage)
	if err != nil {
		return err
	}

	// Construct API URL
	url := fmt.Sprintf(config.Cfg.Telegram.BaseUrl, config.Cfg.Telegram.ApiKey) + config.Cfg.Telegram.EndpointSendMessage

	// Send JSON to API URL
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}

	// Handle HTTP error
	if response.StatusCode != http.StatusOK {
		responseBody := &responseBody{}
		if err = json.NewDecoder(response.Body).Decode(responseBody); err != nil {
			return fmt.Errorf("SendMessage failed with %s (could not decode error response body)", response.Status)
		}
		return fmt.Errorf("SendMessage failed with %d %s", responseBody.ErrorCode, responseBody.Description)
	}

	// No error, we're happy
	return nil
}
