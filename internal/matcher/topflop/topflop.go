package topflop

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/internal/matcher/abstract"
	"github.com/br0-space/bot/internal/telegram"
	"regexp"
	"strings"
)

const identifier = "topflop"

var pattern = regexp.MustCompile(`(?i)^/(top|flop)(@\w+)?($| )`)

var help = []interfaces.MatcherHelpStruct{{
	Command:     `top`,
	Description: `Zeigt eine Liste der am meisten geplusten Begriffe an.`,
}, {
	Command:     `flop`,
	Description: `Zeigt eine Liste der am meisten geminusten Begriffe an.`,
}}

const template = "```\n%s\n```"

type Matcher struct {
	abstract.Matcher
	cfg  interfaces.TopflopMatcherConfigStruct
	repo interfaces.PlusplusRepoInterface
}

func NewMatcher(
	logger interfaces.LoggerInterface,
	config interfaces.TopflopMatcherConfigStruct,
	repo interfaces.PlusplusRepoInterface,
) *Matcher {
	return &Matcher{
		Matcher: abstract.NewMatcher(logger, identifier, pattern, help),
		cfg:     config,
		repo:    repo,
	}
}

func (m *Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) (*[]interfaces.TelegramMessageStruct, error) {
	match := m.GetCommandMatch(messageIn)
	if match == nil {
		return nil, fmt.Errorf("message does not match")
	}

	var records []interfaces.Plusplus
	var err error
	switch match[0] {
	case "top":
		records, err = m.repo.FindTops(messageIn.Chat.ID)
	case "flop":
		records, err = m.repo.FindFlops(messageIn.Chat.ID)
	}
	if err != nil {
		return nil, err
	}

	return reply(records)
}

func reply(records []interfaces.Plusplus) (*[]interfaces.TelegramMessageStruct, error) {
	var lines []string
	for _, record := range records {
		lines = append(lines, fmt.Sprintf(
			"%5d | %s",
			record.Value,
			telegram.EscapeMarkdown(record.Name),
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
