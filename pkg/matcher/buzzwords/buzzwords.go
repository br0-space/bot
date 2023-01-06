package buzzwords

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"github.com/br0-space/bot/pkg/telegram"
	"github.com/mpvl/unique"
	"regexp"
	"strings"
)

const identifier = "buzzwords"

var pattern *regexp.Regexp

var help []interfaces.MatcherHelpStruct

type Matcher struct {
	abstract.Matcher
	repo interfaces.PlusplusRepoInterface
	cfg  Config
}

func NewMatcher(
	logger interfaces.LoggerInterface,
	repo interfaces.PlusplusRepoInterface,
) Matcher {
	var cfg Config
	abstract.LoadMatcherConfig(identifier, &cfg)

	foo := fmt.Sprintf(`(?i)\b(%s)\b`, cfg.GetPattern())
	pattern = regexp.MustCompile(
		foo,
	)

	return Matcher{
		Matcher: abstract.NewMatcher(logger, identifier, pattern, help).WithConfig(&cfg.Config),
		repo:    repo,
		cfg:     cfg,
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) (*[]interfaces.TelegramMessageStruct, error) {
	matches := m.GetInlineMatches(messageIn)
	triggers := m.parseTriggers(matches)

	return m.processTriggers(triggers)
}

func (m Matcher) parseTriggers(matches []string) []string {
	var triggers []string

	for _, match := range matches {
		triggers = append(triggers, m.cfg.GetTrigger(match))
	}

	unique.Sort(unique.StringSlice{P: &triggers})

	return triggers
}

func (m Matcher) processTriggers(triggers []string) (*[]interfaces.TelegramMessageStruct, error) {
	var messages []interfaces.TelegramMessageStruct

	for _, match := range triggers {
		message, err := m.processTrigger(match)
		if err != nil {
			return nil, err
		}
		messages = append(messages, *message)
	}

	return &messages, nil
}

func (m Matcher) processTrigger(trigger string) (*interfaces.TelegramMessageStruct, error) {
	value, err := m.repo.Increment(trigger, 1)
	if err != nil {
		return nil, err
	}

	template, err := m.cfg.GetReply(trigger)
	if err != nil {
		return nil, err
	}

	reply := reply(template, trigger, value)

	return &reply, nil
}

func reply(template string, match string, value int) interfaces.TelegramMessageStruct {
	match = strings.ToUpper(match[:1]) + match[1:]

	return telegram.NewMarkdownMessage(
		fmt.Sprintf(
			template,
			value,
		),
	)
}
