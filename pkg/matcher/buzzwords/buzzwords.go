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

func MakeMatcher(
	logger interfaces.LoggerInterface,
	repo interfaces.PlusplusRepoInterface,
) Matcher {
	var cfg Config
	abstract.LoadMatcherConfig(identifier, &cfg)

	pattern = regexp.MustCompile(fmt.Sprintf(`(?i)\b((%s)([+]{2,}|[-]{2,}|\+-|—)?)`, cfg.GetPattern()))

	return Matcher{
		Matcher: abstract.MakeMatcher(logger, identifier, pattern, help).WithConfig(&cfg.Config),
		repo:    repo,
		cfg:     cfg,
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) ([]interfaces.TelegramMessageStruct, error) {
	matches := m.GetInlineMatches(messageIn)
	triggers := m.parseTriggers(matches)

	return m.makeRepliesFromTriggers(triggers)
}

func (m Matcher) parseTriggers(matches []string) []string {
	var triggers []string

	for _, match := range matches {
		if trigger := m.cfg.GetTrigger(match); trigger != "" {
			triggers = append(triggers, trigger)
		}
	}

	unique.Sort(unique.StringSlice{P: &triggers})

	return triggers
}

func (m Matcher) makeRepliesFromTriggers(triggers []string) ([]interfaces.TelegramMessageStruct, error) {
	var replies []interfaces.TelegramMessageStruct

	for _, match := range triggers {
		triggerReplies, err := m.makeRepliesFromTrigger(match)
		if err != nil {
			return nil, err
		}

		for _, reply := range triggerReplies {
			replies = append(replies, reply)
		}
	}

	return replies, nil
}

func (m Matcher) makeRepliesFromTrigger(trigger string) ([]interfaces.TelegramMessageStruct, error) {
	value, err := m.repo.Increment(trigger, 1)
	if err != nil {
		return nil, err
	}

	template, err := m.cfg.GetReply(trigger)
	if err != nil {
		return nil, err
	}

	reply := makeReply(template, trigger, value)

	return []interfaces.TelegramMessageStruct{reply}, nil
}

func makeReply(template string, match string, value int) interfaces.TelegramMessageStruct {
	match = strings.ToUpper(match[:1]) + match[1:]

	return telegram.MakeMarkdownMessage(
		fmt.Sprintf(
			template,
			value,
		),
	)
}
