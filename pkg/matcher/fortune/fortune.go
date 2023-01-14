package fortune

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"github.com/br0-space/bot/pkg/telegram"
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
	fortuneService interfaces.FortuneServiceInterface
}

func MakeMatcher(
	logger interfaces.LoggerInterface,
	fortune interfaces.FortuneServiceInterface,
) Matcher {
	return Matcher{
		Matcher:        abstract.MakeMatcher(logger, identifier, pattern, help),
		fortuneService: fortune,
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) ([]interfaces.TelegramMessageStruct, error) {
	match := m.GetCommandMatch(messageIn)
	if match == nil {
		return nil, fmt.Errorf("message does not match")
	}

	switch strings.TrimSpace(match[3]) {
	case "list":
		return m.makeListReplies()
	case "":
		return m.makeRandomReplies()
	default:
		return m.makeFromFileReplies(strings.TrimSpace(match[3]))
	}
}

func (m Matcher) makeListReplies() ([]interfaces.TelegramMessageStruct, error) {
	text := fmt.Sprintf(
		templates.list,
		strings.Join(m.fortuneService.GetList(), "\n"),
	)

	return []interfaces.TelegramMessageStruct{
		telegram.MakeMarkdownMessage(text),
	}, nil
}

func (m Matcher) makeRandomReplies() ([]interfaces.TelegramMessageStruct, error) {
	fortune, err := m.fortuneService.GetRandomFortune()
	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf(
		templates.random,
		fortune.ToMarkdown(),
		fortune.File(),
	)

	return []interfaces.TelegramMessageStruct{
		telegram.MakeMarkdownMessage(text),
	}, nil
}

func (m Matcher) makeFromFileReplies(file string) ([]interfaces.TelegramMessageStruct, error) {
	if !m.fortuneService.Exists(file) {
		return nil, fmt.Errorf(`fortune file "%s" does not exist`, file)
	}

	fortune, err := m.fortuneService.GetFortune(file)
	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf(
		templates.fromFile,
		fortune.ToMarkdown(),
	)

	return []interfaces.TelegramMessageStruct{
		telegram.MakeMarkdownMessage(text),
	}, nil
}
