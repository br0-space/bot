package stats

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"github.com/br0-space/bot/pkg/telegram"
	"regexp"
	"strings"
)

const identifier = "stats"

var pattern = regexp.MustCompile(`(?i)^/(stats)(@\w+)?($| )`)

var help = []interfaces.MatcherHelpStruct{{
	Description: `Zeigt eine Liste der dem Bot bekannten User an.`,
}}

const template = "```\n%s\n```"

type Matcher struct {
	abstract.Matcher
	repo interfaces.StatsRepoInterface
}

func NewMatcher(
	logger interfaces.LoggerInterface,
	repo interfaces.StatsRepoInterface,
) Matcher {
	return Matcher{
		Matcher: abstract.NewMatcher(logger, identifier, pattern, help),
		repo:    repo,
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) (*[]interfaces.TelegramMessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, fmt.Errorf("message does not match")
	}

	users, err := m.repo.GetTopUsers()
	if err != nil {
		return nil, err
	}

	return reply(users)
}

func reply(users []interfaces.StatsUserStruct) (*[]interfaces.TelegramMessageStruct, error) {
	var lines []string
	for _, user := range users {
		lines = append(lines, fmt.Sprintf(
			"%6d | %s",
			user.Posts,
			telegram.EscapeMarkdown(user.Username),
		))
	}

	text := fmt.Sprintf(
		template,
		strings.Join(lines, "\n"),
	)

	return &[]interfaces.TelegramMessageStruct{
		telegram.NewMarkdownMessage(text),
	}, nil
}