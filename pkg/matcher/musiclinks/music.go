package musiclinks

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"github.com/br0-space/bot/pkg/songlink"
	"github.com/br0-space/bot/pkg/telegram"
	"regexp"
)

const identifier = "musiclinks"

var pattern = regexp.MustCompile(`(https?://open.spotify.com/(album|track)/.+?|https?://music.apple.com/[a-z]{2}/album/.+?)(\s|$)`)

var help []interfaces.MatcherHelpStruct

type Matcher struct {
	abstract.Matcher
}

func NewMatcher(logger interfaces.LoggerInterface) Matcher {
	return Matcher{
		Matcher: abstract.NewMatcher(logger, identifier, pattern, help),
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) (*[]interfaces.TelegramMessageStruct, error) {
	res := make([]interfaces.TelegramMessageStruct, 0)

	matches := m.GetInlineMatches(messageIn)
	for _, match := range matches {
		songlinkEntry, err := songlink.GetSonglinkEntry(match)
		if err != nil {
			return nil, err
		}

		res = append(res, reply(*songlinkEntry))
	}

	return &res, nil
}

func reply(songlinkEntry songlink.Entry) interfaces.TelegramMessageStruct {
	text := fmt.Sprintf(
		"*%s*\n*%s* · %s\n\n",
		telegram.EscapeMarkdown(songlinkEntry.Title),
		telegram.EscapeMarkdown(songlinkEntry.Artist),
		songlinkEntry.Type.Natural(),
	)

	for i := range songlinkEntry.Links {
		if songlinkEntry.Links[i].Platform == songlink.Songlink {
			continue
		}

		text += fmt.Sprintf(
			"🎧 [%s](%s)\n\n",
			songlinkEntry.Links[i].Platform.Natural(),
			songlinkEntry.Links[i].URL,
		)
	}

	text = text + fmt.Sprintf(
		"🔗 [%s](%s)",
		songlink.Songlink.Natural(),
		songlinkEntry.Links[0].URL,
	)

	res := telegram.NewMarkdownMessage(text)
	res.DisableWebPagePreview = true

	return res
}
