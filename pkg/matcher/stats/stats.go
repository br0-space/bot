package stats

import (
	"fmt"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
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
	repo interfaces.UserStatsRepoInterface
}

func MakeMatcher(
	repo interfaces.UserStatsRepoInterface,
) Matcher {
	return Matcher{
		Matcher: abstract.MakeMatcher(identifier, pattern, help),
		repo:    repo,
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, fmt.Errorf("message does not match")
	}

	users, err := m.repo.GetTopUsers()
	if err != nil {
		return nil, err
	}

	return makeReplies(users)
}

func makeReplies(users []interfaces.StatsUserStruct) ([]telegramclient.MessageStruct, error) {
	var lines []string
	for _, user := range users {
		lines = append(lines, fmt.Sprintf(
			"%6d | %s",
			user.Posts,
			telegramclient.EscapeMarkdown(user.Username),
		))
	}

	text := fmt.Sprintf(
		template,
		strings.Join(lines, "\n"),
	)

	return []telegramclient.MessageStruct{
		telegramclient.MarkdownMessage(text),
	}, nil
}
