package abstract

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

type Matcher struct {
	logger     interfaces.LoggerInterface
	identifier string
	regexp     *regexp.Regexp
	help       []interfaces.MatcherHelpStruct
	cfg        *Config
}

type Config struct {
	enabled *bool
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

func (m Matcher) WithConfig(cfg *Config) Matcher {
	m.cfg = cfg

	return m
}

func (m Matcher) IsEnabled() bool {
	if m.cfg == nil || m.cfg.enabled == nil {
		return true
	}

	return *m.cfg.enabled
}

func (m Matcher) GetLogger() interfaces.LoggerInterface {
	return m.logger
}

func (m Matcher) GetIdentifier() string {
	return m.identifier
}

func (m Matcher) GetHelp() []interfaces.MatcherHelpStruct {
	return m.help
}

func (m Matcher) DoesMatch(messageIn interfaces.TelegramWebhookMessageStruct) bool {
	return m.regexp.MatchString(messageIn.Text)
}

func (m Matcher) GetCommandMatch(messageIn interfaces.TelegramWebhookMessageStruct) []string {
	match := m.regexp.FindStringSubmatch(messageIn.Text)
	if match == nil {
		return nil
	}
	if len(match) > 0 {
		return match[1:]
	}

	return match
}

func (m Matcher) GetInlineMatches(messageIn interfaces.TelegramWebhookMessageStruct) []string {
	matches := m.regexp.FindAllString(messageIn.Text, -1)
	if matches == nil {
		return []string{}
	}

	for i, match := range matches {
		matches[i] = strings.TrimSpace(match)
	}

	return matches
}

func (m Matcher) HandleError(_ interfaces.TelegramWebhookMessageStruct, identifier string, err error) {
	m.logger.Error(identifier, err.Error())
}

func LoadMatcherConfig(identifier string, cfg interface{}) {
	v := viper.New()

	v.SetConfigFile(fmt.Sprintf("config/%s.yaml", identifier))
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		panic(err)
	}
}
