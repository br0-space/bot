package fortune

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/internal/fortune"
	"github.com/br0-space/bot/internal/matcher/abstract"
	"github.com/br0-space/bot/internal/telegram"
	"regexp"
	"strings"
)

const identifier = "fortune"

var pattern = regexp.MustCompile(`(?i)^/(fortune)(@\w+)?($| )(.+)?$`)

var help = []interfaces.MatcherHelpStruct{{
	Command:     `fortune`,
	Description: `Zeigt ein Fortune Cookie an.`,
	Usage:       `/fortune (list|<optional: File>)`,
	Example:     `/fortune wisdom`,
}}

var templates = struct {
	list     string
	random   string
	fromFile string
}{
	list:     "*Available Fortune Cookie Files*\n\n%s",
	random:   "%s\n\n_\\[from `%s`\\]_",
	fromFile: "%s",
}

type Matcher struct {
	abstract.Matcher
}

func NewMatcher(logger interfaces.LoggerInterface) Matcher {
	return Matcher{
		Matcher: abstract.NewMatcher(logger, identifier, pattern, help),
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) (*[]interfaces.TelegramMessageStruct, error) {
	match := m.GetCommandMatch(messageIn)
	if match == nil {
		return nil, fmt.Errorf("message does not match")
	}

	switch strings.TrimSpace(match[3]) {
	case "list":
		return replyList()
	case "":
		return replyRandom()
	default:
		return replyFromFile(strings.TrimSpace(match[3]))
	}
}

func replyList() (*[]interfaces.TelegramMessageStruct, error) {
	text := fmt.Sprintf(
		templates.list,
		strings.Join(fortune.GetList(), "\n"),
	)

	return &[]interfaces.TelegramMessageStruct{
		telegram.NewMarkdownMessage(text),
	}, nil
}

func replyRandom() (*[]interfaces.TelegramMessageStruct, error) {
	res, err := fortune.GetRandomFortune()
	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf(
		templates.random,
		res.ToMarkdown(),
		res.GetFile(),
	)

	return &[]interfaces.TelegramMessageStruct{
		telegram.NewMarkdownMessage(text),
	}, nil
}

func replyFromFile(file string) (*[]interfaces.TelegramMessageStruct, error) {
	if !fortune.Exists(file) {
		return nil, fmt.Errorf(`fortune file "%s" does not exist`, file)
	}

	res, err := fortune.GetFortune(file)
	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf(
		templates.fromFile,
		res.ToMarkdown(),
	)

	return &[]interfaces.TelegramMessageStruct{
		telegram.NewMarkdownMessage(text),
	}, nil
}
