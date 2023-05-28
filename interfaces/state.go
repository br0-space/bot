package interfaces

import (
	telegramclient "github.com/br0-space/bot-telegramclient"
	"time"
)

type StateServiceInterface interface {
	ProcessMessage(messageIn telegramclient.WebhookMessageStruct)
	GetLastPost(userID int64) *time.Time
}
