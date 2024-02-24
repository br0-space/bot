package stats

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
)

const identifier = "stats"

var pattern = regexp.MustCompile(`(?i)^/(stats)(@\w+)?($| )`)

var help = []matcher.HelpStruct{{
	Command:     "",
	Description: `Zeigt eine Liste der dem Bot bekannten User an.`,
	Usage:       "",
	Example:     "",
}}

const template = "```\n%s\n```"

type Matcher struct {
	matcher.Matcher
	repo interfaces.UserStatsRepoInterface
}

func MakeMatcher(
	repo interfaces.UserStatsRepoInterface,
) Matcher {
	return Matcher{
		Matcher: matcher.MakeMatcher(identifier, pattern, help),
		repo:    repo,
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, errors.New("message does not match")
	}

	users, err := m.repo.GetTopUsers()
	if err != nil {
		return nil, err
	}

	return makeReplies(users)
}

func makeReplies(users []interfaces.StatsUserStruct) ([]telegramclient.MessageStruct, error) {
	lines := make([]string, 0, len(users))
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
