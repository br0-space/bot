package interfaces

import (
	"time"

	telegramclient "github.com/br0-space/bot-telegramclient"
)

type StateServiceInterface interface {
	ProcessMessage(messageIn telegramclient.WebhookMessageStruct)
	GetLastPost(userID int64) *time.Time
}
