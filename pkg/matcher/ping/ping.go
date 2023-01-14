package ping

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"github.com/br0-space/bot/pkg/telegram"
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

func MakeMatcher(logger interfaces.LoggerInterface) Matcher {
	return Matcher{
		Matcher: abstract.MakeMatcher(logger, identifier, pattern, help),
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) ([]interfaces.TelegramMessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, fmt.Errorf("message does not match")
	}

	return []interfaces.TelegramMessageStruct{
		telegram.MakeReply(template, messageIn.ID),
	}, nil
}
