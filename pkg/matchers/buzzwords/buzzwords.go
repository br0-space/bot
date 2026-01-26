package buzzwords

import (
	"fmt"
	"regexp"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
	"github.com/mpvl/unique"
)

const identifier = "buzzwords"

var pattern *regexp.Regexp

var help []matcher.HelpStruct

type Matcher struct {
	matcher.WithCustomConfigType[Config]

	repo interfaces.PlusplusRepoInterface
}

func MakeMatcher(
	repo interfaces.PlusplusRepoInterface,
) Matcher {
	cfgs, err := matcher.LoadMatcherConfig[Config](identifier)
	if err != nil {
		panic(fmt.Sprintf("failed to load matcher config for %s: %v", identifier, err))
	}

	cfg := cfgs[0]

	pattern = regexp.MustCompile(fmt.Sprintf(`(?i)\b((%s)([+]{2,}|[-]{2,}|\+-|—)?)`, cfg.GetPattern()))

	return Matcher{
		WithCustomConfigType: matcher.MakeMatcherWithCustomConfigType(identifier, pattern, help, cfg),
		repo:                 repo,
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	matches := m.InlineMatches(messageIn)
	triggers := m.parseTriggers(matches)

	return m.makeRepliesFromTriggers(triggers)
}

func (m Matcher) parseTriggers(matches []string) []string {
	var triggers []string

	for _, match := range matches {
		if trigger := m.Config().GetTrigger(match); trigger != "" {
			triggers = append(triggers, trigger)
		}
	}

	unique.Sort(unique.StringSlice{P: &triggers})

	return triggers
}

func (m Matcher) makeRepliesFromTriggers(triggers []string) ([]telegramclient.MessageStruct, error) {
	var replies []telegramclient.MessageStruct

	for _, match := range triggers {
		triggerReplies, err := m.makeRepliesFromTrigger(match)
		if err != nil {
			return nil, err
		}

		replies = append(replies, triggerReplies...)
	}

	return replies, nil
}

func (m Matcher) makeRepliesFromTrigger(trigger string) ([]telegramclient.MessageStruct, error) {
	value, err := m.repo.Increment(trigger, 1)
	if err != nil {
		return nil, err
	}

	template, err := m.Config().GetReply(trigger)
	if err != nil {
		return nil, err
	}

	reply := makeReply(template, value)

	return []telegramclient.MessageStruct{reply}, nil
}

func makeReply(template string, value int) telegramclient.MessageStruct {
	return telegramclient.MarkdownMessage(
		fmt.Sprintf(
			template,
			value,
		),
	)
}
