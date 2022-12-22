package atall

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/internal/matcher/abstract"
	"github.com/br0-space/bot/internal/telegram"
	"regexp"
	"strings"
)

const identifier = "atall"

var pattern = regexp.MustCompile(`(^|\s)@alle?(\s|$)`)

var help []interfaces.MatcherHelpStruct

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
	matches := m.GetInlineMatches(messageIn)
	if len(matches) == 0 {
		return nil, nil
	}

	users, err := m.repo.GetKnownUsers()
	if err != nil {
		return nil, err
	}

	return reply(messageIn.TextOrCaption(), users)
}

func reply(text string, users []interfaces.StatsUserStruct) (*[]interfaces.TelegramMessageStruct, error) {
	text = strings.ReplaceAll(text, "@alle", "")
	text = strings.ReplaceAll(text, "@all", "")

	for _, user := range users {
		text += fmt.Sprintf(
			" [%s](tg://user?id=%d)",
			telegram.EscapeMarkdown(user.Username),
			user.ID,
		)
	}

	return &[]interfaces.TelegramMessageStruct{
		telegram.NewMarkdownMessage(text),
	}, nil
}
