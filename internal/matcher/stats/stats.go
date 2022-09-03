package stats

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/internal/matcher/abstract"
	"github.com/br0-space/bot/internal/telegram"
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
	cfg  interfaces.StatsMatcherConfigStruct
	repo interfaces.StatsRepoInterface
}

func NewMatcher(
	logger interfaces.LoggerInterface,
	config interfaces.StatsMatcherConfigStruct,
	repo interfaces.StatsRepoInterface,
) *Matcher {
	return &Matcher{
		Matcher: abstract.NewMatcher(logger, identifier, pattern, help),
		cfg:     config,
		repo:    repo,
	}
}

func (m *Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) (*[]interfaces.TelegramMessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, fmt.Errorf("message does not match")
	}

	users, err := m.repo.GetTopUsers(messageIn.Chat.ID)
	if err != nil {
		return nil, err
	}

	return reply(users)
}

func reply(users []interfaces.StatsUserStruct) (*[]interfaces.TelegramMessageStruct, error) {
	var lines []string
	for _, user := range users {
		lines = append(lines, fmt.Sprintf("%6d | %s", user.Posts, user.Username))
	}

	text := fmt.Sprintf(
		template,
		strings.Join(lines, "\n"),
	)

	return &[]interfaces.TelegramMessageStruct{
		telegram.NewMarkdownMessage(text),
	}, nil
}
