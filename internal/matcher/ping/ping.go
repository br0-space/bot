package ping

import (
	"regexp"

	"github.com/br0-space/bot/container"
	"github.com/br0-space/bot/internal/logger"
	"github.com/br0-space/bot/internal/matcher"
	"github.com/br0-space/bot/internal/telegram"
	"github.com/br0-space/bot/internal/telegram/webhook"
	"github.com/segmentio/stats/v4"
)

type Matcher struct {
	matcher.Matcher
	log      logger.Interface
	telegram telegram.Interface
}

func MakeMatcher() Matcher {
	return Matcher{
		log:      container.ProvideLoggerService(),
		telegram: container.ProvideTelegramService(),
	}
}

func (m Matcher) Identifier() string {
	return "ping"
}

func (m Matcher) ProcessMessage(messageIn webhook.Message) error {
	if doesMatch := m.doesMatch(messageIn.Text); !doesMatch {
		return nil
	}

	stats.Incr("ping")

	return m.sendResponse(messageIn)
}

func (m Matcher) doesMatch(text string) bool {
	match, _ := regexp.MatchString(`^/ping(@|\s|$)`, text)
	return match
}

func (m Matcher) sendResponse(messageIn webhook.Message) error {
	messageOut := telegram.Message{
		Text:             "pong",
		ReplyToMessageID: messageIn.ID,
	}
	return m.telegram.SendMessage(messageIn.Chat.ID, messageOut)
}
