package abstract

import (
	"github.com/br0-space/bot/interfaces"
	"regexp"
	"strings"
)

type Matcher struct {
	logger     interfaces.LoggerInterface
	identifier string
	regexp     *regexp.Regexp
	help       []interfaces.MatcherHelpStruct
}

func NewMatcher(
	logger interfaces.LoggerInterface,
	identifier string,
	pattern *regexp.Regexp,
	help []interfaces.MatcherHelpStruct,
) Matcher {
	return Matcher{
		logger:     logger,
		identifier: identifier,
		regexp:     pattern,
		help:       help,
	}
}

func (m *Matcher) GetLogger() interfaces.LoggerInterface {
	return m.logger
}

func (m *Matcher) GetIdentifier() string {
	return m.identifier
}

func (m *Matcher) GetHelp() []interfaces.MatcherHelpStruct {
	return m.help
}

func (m *Matcher) DoesMatch(messageIn interfaces.TelegramWebhookMessageStruct) bool {
	return m.regexp.MatchString(messageIn.Text)
}

func (m *Matcher) GetCommandMatch(messageIn interfaces.TelegramWebhookMessageStruct) []string {
	match := m.regexp.FindStringSubmatch(messageIn.Text)
	if match == nil {
		return nil
	}
	if len(match) > 0 {
		return match[1:]
	}

	return match
}

func (m *Matcher) GetInlineMatches(messageIn interfaces.TelegramWebhookMessageStruct) []string {
	matches := m.regexp.FindAllString(messageIn.Text, -1)
	if matches == nil {
		return []string{}
	}

	for i, match := range matches {
		matches[i] = strings.TrimSpace(match)
	}

	return matches
}

func (m *Matcher) HandleError(messageIn interfaces.TelegramWebhookMessageStruct, identifier string, err error) {
	m.logger.Error(identifier, err.Error())
}
