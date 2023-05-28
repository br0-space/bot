package musiclinks

import (
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"regexp"
)

const identifier = "musiclinks"

var pattern = regexp.MustCompile(`(https?://open.spotify.com/(album|track)/.+?|https?://music.apple.com/[a-z]{2}/album/.+?)(\s|$)`)

var help []interfaces.MatcherHelpStruct

type Matcher struct {
	abstract.Matcher
	songlinkService interfaces.SonglinkServiceInterface
}

func MakeMatcher(
	songlinkService interfaces.SonglinkServiceInterface,
) Matcher {
	return Matcher{
		Matcher:         abstract.MakeMatcher(identifier, pattern, help),
		songlinkService: songlinkService,
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	matches := m.GetInlineMatches(messageIn)

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
