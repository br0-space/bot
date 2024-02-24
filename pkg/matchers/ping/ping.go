package ping

import (
	"errors"
	"regexp"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
)

const identifier = "ping"

var pattern = regexp.MustCompile(`(?i)^/(ping)(@\w+)?($| )`)

var help = []matcher.HelpStruct{{
	Command:     "",
	Description: `Antwortet mit "pong"`,
	Usage:       "",
	Example:     "",
}}

const template = `pong`

type Matcher struct {
	matcher.Matcher
}

func MakeMatcher() Matcher {
	return Matcher{
		Matcher: matcher.MakeMatcher(identifier, pattern, help),
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, errors.New("message does not match")
	}

	return []telegramclient.MessageStruct{
		telegramclient.Reply(template, messageIn.ID),
	}, nil
}
