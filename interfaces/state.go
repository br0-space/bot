package interfaces

import "time"

type StateServiceInterface interface {
	ProcessMessage(messageIn TelegramWebhookMessageStruct)
	GetLastPost(userID int64) *time.Time
}
