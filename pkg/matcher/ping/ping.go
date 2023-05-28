package ping

import (
	"fmt"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
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

func MakeMatcher() Matcher {
	return Matcher{
		Matcher: abstract.MakeMatcher(identifier, pattern, help),
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, fmt.Errorf("message does not match")
	}

	return []telegramclient.MessageStruct{
		telegramclient.Reply(template, messageIn.ID),
	}, nil
}
