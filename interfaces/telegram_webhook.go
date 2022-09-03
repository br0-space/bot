package interfaces

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type TelegramWebhookHandlerInterface interface {
	InitMatchers()
	ServeHTTP(res http.ResponseWriter, req *http.Request)
}

type TelegramWebhookToolsInterface interface {
	SetWebhookURL()
}

// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update

type TelegramWebhookBodyStruct struct {
	Message TelegramWebhookMessageStruct `json:"message"`
}

type TelegramWebhookMessageStruct struct {
	ID      int64                               `json:"message_id"`
	From    TelegramWebhookMessageUserStruct    `json:"from"`
	Chat    TelegramWebhookMessageChatStruct    `json:"chat"`
	Text    string                              `json:"text"`
	Date    int64                               `json:"date"`
	Photo   []TelegramWebhookMessagePhotoStruct `json:"photo"`
	Caption string                              `json:"caption"`
}

type TelegramWebhookMessageUserStruct struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type TelegramWebhookMessageChatStruct struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Username string `json:"username"`
}

type TelegramWebhookMessagePhotoStruct struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

// Helper functions

func (m TelegramWebhookMessageStruct) TextOrCaption() string {
	if len(m.Text) > 0 {
		return m.Text
	}

	if len(m.Caption) > 0 {
		return m.Caption
	}

	return ""
}

func (m TelegramWebhookMessageStruct) WordCount() int {
	// Match non-space character sequences.
	re := regexp.MustCompile(`[\S]+`)

	// Find all matches and return count.
	results := re.FindAllString(m.TextOrCaption(), -1)
	return len(results)
}

func (u TelegramWebhookMessageUserStruct) UsernameOrName() string {
	if len(u.Username) > 0 {
		return "@" + u.Username
	}

	return strings.Trim(fmt.Sprintf("%s %s", u.FirstName, u.LastName), " ")
}

// Test functions

func NewTestTelegramWebhookMessage(text string) TelegramWebhookMessageStruct {
	return TelegramWebhookMessageStruct{
		ID:   123,
		From: NewTestTelegramWebhookMessageUser(false),
		Chat: NewTestTelegramWebhookMessageChat(),
		Text: text,
	}
}

func NewTestTelegramWebhookMessageUser(isBot bool) TelegramWebhookMessageUserStruct {
	return TelegramWebhookMessageUserStruct{
		ID:       456,
		IsBot:    isBot,
		Username: "Foobar",
	}
}

func NewTestTelegramWebhookMessageChat() TelegramWebhookMessageChatStruct {
	return TelegramWebhookMessageChatStruct{
		ID: 789,
	}
}
