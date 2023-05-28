package xkcd

import (
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"regexp"
	"strconv"
	"strings"
)

const identifier = "xkcd"

var pattern = regexp.MustCompile(`(?i)^/(xkcd)(@\w+)?($| )(.+)?$`)

var help = []interfaces.MatcherHelpStruct{{
	Command:     `xkcd`,
	Description: `Zeigt einen xkcd Comic an.`,
	Usage:       `/xkcd (latest|<optional: Comic ID>)`,
	Example:     `/xkcd 1234`,
}}

type Matcher struct {
	abstract.Matcher
	xkcdService interfaces.XkcdServiceInterface
}

func MakeMatcher(
	xkcd interfaces.XkcdServiceInterface,
) Matcher {
	return Matcher{
		Matcher:     abstract.MakeMatcher(identifier, pattern, help),
		xkcdService: xkcd,
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	match := m.GetCommandMatch(messageIn)
	if match == nil {
		return nil, nil
	}

	subCommand := strings.TrimSpace(match[3])

	switch {
	case subCommand == "latest":
		return m.makeLatestReplies()
	case regexp.MustCompile(`^\d+$`).MatchString(subCommand):
		if id, err := strconv.Atoi(subCommand); err != nil {
			return nil, err
		} else {
			return m.makeFromIDReplies(id)
		}
	default:
		return m.makeRandomReplies()
	}
}

func (m Matcher) makeLatestReplies() ([]telegramclient.MessageStruct, error) {
	comic, err := m.xkcdService.Latest()
	if err != nil {
		return nil, err
	}

	return m.makeReplies(comic), nil
}

func (m Matcher) makeFromIDReplies(id int) ([]telegramclient.MessageStruct, error) {
	comic, err := m.xkcdService.Comic(id)
	if err != nil {
		return nil, err
	}

	return m.makeReplies(comic), nil
}

func (m Matcher) makeRandomReplies() ([]telegramclient.MessageStruct, error) {
	comic, err := m.xkcdService.Random()
	if err != nil {
		return nil, err
	}

	return m.makeReplies(comic), nil
}

func (m Matcher) makeReplies(comic interfaces.XkcdComicInterface) []telegramclient.MessageStruct {
	return []telegramclient.MessageStruct{
		telegramclient.MarkdownPhoto(comic.ImageURL(), comic.ToMarkdown()),
	}
}
