package atall

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"github.com/br0-space/bot/pkg/telegram"
	"regexp"
	"strings"
)

const identifier = "atall"

var pattern = regexp.MustCompile(`(^|\s)@alle?(\s|$)`)

var help []interfaces.MatcherHelpStruct

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

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) ([]interfaces.TelegramMessageStruct, error) {
	matches := m.GetInlineMatches(messageIn)
	if len(matches) == 0 {
		return nil, nil
	}

	users, err := m.repo.GetKnownUsers()
	if err != nil {
		return nil, err
	}

	return makeReplies(messageIn.TextOrCaption(), users)
}

func makeReplies(text string, users []interfaces.StatsUserStruct) ([]interfaces.TelegramMessageStruct, error) {
	text = strings.ReplaceAll(text, "@alle", "")
	text = strings.ReplaceAll(text, "@all", "")

	for _, user := range users {
		text += fmt.Sprintf(
			" [%s](tg://user?id=%d)",
			telegram.EscapeMarkdown(user.Username),
			user.ID,
		)
	}

	return []interfaces.TelegramMessageStruct{
		telegram.MakeMarkdownMessage(text),
	}, nil
}
