package topflop

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
)

const (
	identifier   = "topflop"
	defaultLimit = 10
)

var pattern = regexp.MustCompile(`(?i)^/(top|flop)(@\w+)?($| )(\d+)?`)

var help = []matcher.HelpStruct{{
	Command:     `top`,
	Description: `Zeigt eine Liste der am meisten geplusten Begriffe an.`,
	Usage:       `/top <optional: Anzahl der Einträge>`,
	Example:     `/top 10`,
}, {
	Command:     `flop`,
	Description: `Zeigt eine Liste der am meisten geminusten Begriffe an.`,
	Usage:       `/flop <optional: Anzahl der Einträge>`,
	Example:     `/flop 10`,
}}

const template = "```\n%s\n```"

type Matcher struct {
	matcher.Matcher
	repo interfaces.PlusplusRepoInterface
}

func MakeMatcher(
	repo interfaces.PlusplusRepoInterface,
) Matcher {
	return Matcher{
		Matcher: matcher.MakeMatcher(identifier, pattern, help),
		repo:    repo,
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	match := m.CommandMatch(messageIn)
	if match == nil {
		return nil, errors.New("message does not match")
	}

	cmd := match[0]
	limit := defaultLimit

	if match[3] != "" {
		res, err := strconv.ParseInt(match[3], 10, 0)
		if err != nil {
			return nil, err
		}

		limit = int(res)
	}

	var records []interfaces.Plusplus

	var err error

	switch cmd {
	case "top":
		records, err = m.repo.FindTops(limit)
	case "flop":
		records, err = m.repo.FindFlops(limit)
	}

	if err != nil {
		return nil, err
	}

	return makeReplies(records, messageIn.ID)
}

func makeReplies(records []interfaces.Plusplus, messageID int64) ([]telegramclient.MessageStruct, error) {
	lines := make([]string, 0, len(records))
	for _, record := range records {
		lines = append(lines, fmt.Sprintf(
			"%5d | %s",
			record.Value,
			telegramclient.EscapeMarkdown(record.Name),
		))
	}

	text := fmt.Sprintf(
		template,
		strings.Join(lines, "\n"),
	)

	return []telegramclient.MessageStruct{
		telegramclient.MarkdownReply(text, messageID),
	}, nil
}
