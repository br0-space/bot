package atall

import (
	"fmt"
	"regexp"
	"strings"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
)

const identifier = "atall"

var pattern = regexp.MustCompile(`(^|\s)@alle?(\s|$)`)

var help []matcher.HelpStruct

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
	matches := m.InlineMatches(messageIn)
	if len(matches) == 0 {
		return nil, nil
	}

	users, err := m.repo.GetKnownUsers()
	if err != nil {
		return nil, err
	}

	return makeReplies(messageIn.TextOrCaption(), users)
}

func makeReplies(text string, users []interfaces.StatsUserStruct) ([]telegramclient.MessageStruct, error) {
	text = strings.ReplaceAll(text, "@alle", "")
	text = strings.ReplaceAll(text, "@all", "")

	for _, user := range users {
		text += fmt.Sprintf(
			" [%s](tg://user?id=%d)",
			telegramclient.EscapeMarkdown(user.Username),
			user.ID,
		)
	}

	return []telegramclient.MessageStruct{
		telegramclient.MarkdownMessage(text),
	}, nil
}
