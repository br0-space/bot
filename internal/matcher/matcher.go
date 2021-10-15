package matcher

import (
	"fmt"

	"github.com/br0-space/bot/container"
	"github.com/br0-space/bot/internal/logger"
	"github.com/br0-space/bot/internal/telegram"
	"github.com/br0-space/bot/internal/telegram/webhook"
)

// Every matcher must implement this interface
type Interface interface {
	Identifier() string
	ProcessMessage(messageIn webhook.Message) error
	HandleError(messageIn webhook.Message, err error)
}

// Create a struct each matcher is inherited from
type Matcher struct{
	log logger.Interface
	telegram telegram.Interface
}

func MakeMatcher() Matcher {
	return Matcher{
		log: container.ProvideLoggerService(),
		telegram: container.ProvideTelegramService(),
	}
}

func (m Matcher) Identifier() string {
	return "matcher"
}

// Handle an error happening during matcher execution
func (m Matcher) HandleError(messageIn webhook.Message, err error) {
	// Print error information to logger
	m.log.Error(m.Identifier(), err.Error())

	// Send error notification to telegram
	err = m.telegram.SendMessage(
		messageIn.Chat.ID,
		telegram.Message{
			Text:             fmt.Sprintf("⚠️ %s %s", m.Identifier(), err.Error()),
			ReplyToMessageID: messageIn.ID,
		},
	)
	if err != nil {
		m.log.Warning(m.Identifier(), err.Error())
	}
}