package musiclinks

import (
	"regexp"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
)

const identifier = "musiclinks"

var pattern = regexp.MustCompile(`(https?://open.spotify.com/(album|track)/.+?|https?://music.apple.com/[a-z]{2}/album/.+?)(\s|$)`)

var help []matcher.HelpStruct

type Matcher struct {
	matcher.Matcher
	songlinkService interfaces.SonglinkServiceInterface
}

func MakeMatcher(
	songlinkService interfaces.SonglinkServiceInterface,
) Matcher {
	return Matcher{
		Matcher:         matcher.MakeMatcher(identifier, pattern, help),
		songlinkService: songlinkService,
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	matches := m.InlineMatches(messageIn)

	res := make([]telegramclient.MessageStruct, 0)

	for _, match := range matches {
		songlinkEntry, err := m.songlinkService.GetEntryForUrl(match)
		if err != nil {
			return nil, err
		}

		res = append(res, makeReply(songlinkEntry))
	}

	return res, nil
}

func makeReply(songlinkEntry interfaces.SonglinkEntryInterface) telegramclient.MessageStruct {
	res := telegramclient.MarkdownMessage(songlinkEntry.ToMarkdown())
	res.DisableWebPagePreview = true

	return res
}
