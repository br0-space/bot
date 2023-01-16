package xkcd

import (
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"github.com/br0-space/bot/pkg/telegram"
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
	logger interfaces.LoggerInterface,
	xkcd interfaces.XkcdServiceInterface,
) Matcher {
	return Matcher{
		Matcher:     abstract.MakeMatcher(logger, identifier, pattern, help),
		xkcdService: xkcd,
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) ([]interfaces.TelegramMessageStruct, error) {
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

func (m Matcher) makeLatestReplies() ([]interfaces.TelegramMessageStruct, error) {
	comic, err := m.xkcdService.Latest()
	if err != nil {
		return nil, err
	}

	return m.makeReplies(comic), nil
}

func (m Matcher) makeFromIDReplies(id int) ([]interfaces.TelegramMessageStruct, error) {
	comic, err := m.xkcdService.Comic(id)
	if err != nil {
		return nil, err
	}

	return m.makeReplies(comic), nil
}

func (m Matcher) makeRandomReplies() ([]interfaces.TelegramMessageStruct, error) {
	comic, err := m.xkcdService.Random()
	if err != nil {
		return nil, err
	}

	return m.makeReplies(comic), nil
}

func (m Matcher) makeReplies(comic interfaces.XkcdComicInterface) []interfaces.TelegramMessageStruct {
	return []interfaces.TelegramMessageStruct{
		telegram.MakeMarkdownPhoto(comic.ImageURL(), comic.ToMarkdown()),
	}
}
