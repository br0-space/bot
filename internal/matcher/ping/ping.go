package ping

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/internal/matcher/abstract"
	"github.com/br0-space/bot/internal/telegram"
	"regexp"
)

const identifier = "ping"

var pattern = regexp.MustCompile(`(?i)^/(ping)(@\w+)?($| )`)

var help = []interfaces.MatcherHelpStruct{{
	Description: `Antwortet mit "pong"`,
}}

const template = `pong`

type Matcher struct {
	abstract.Matcher
}

func NewMatcher(logger interfaces.LoggerInterface) Matcher {
	return Matcher{
		Matcher: abstract.NewMatcher(logger, identifier, pattern, help),
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) (*[]interfaces.TelegramMessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, fmt.Errorf("message does not match")
	}

	return &[]interfaces.TelegramMessageStruct{
		telegram.NewReply(template, messageIn.ID),
	}, nil
}
