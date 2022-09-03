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
	cfg interfaces.PingMatcherConfigStruct
}

func NewMatcher(logger interfaces.LoggerInterface, config interfaces.PingMatcherConfigStruct) *Matcher {
	return &Matcher{
		Matcher: abstract.NewMatcher(logger, identifier, pattern, help),
		cfg:     config,
	}
}

func (m *Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) (*[]interfaces.TelegramMessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, fmt.Errorf("message does not match")
	}

	return &[]interfaces.TelegramMessageStruct{
		telegram.NewReply(template, messageIn.ID),
	}, nil
}
