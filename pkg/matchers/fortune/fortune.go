package fortune

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
)

const identifier = "fortune"

var pattern = regexp.MustCompile(`(?i)^/(fortune)(@\w+)?($| )(.+)?$`)

var help = []matcher.HelpStruct{{
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
	matcher.Matcher
	fortuneService interfaces.FortuneServiceInterface
}

func MakeMatcher(
	fortune interfaces.FortuneServiceInterface,
) Matcher {
	return Matcher{
		Matcher:        matcher.MakeMatcher(identifier, pattern, help),
		fortuneService: fortune,
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	match := m.CommandMatch(messageIn)
	if match == nil {
		return nil, errors.New("message does not match")
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

func (m Matcher) makeListReplies() ([]telegramclient.MessageStruct, error) {
	text := fmt.Sprintf(
		templates.list,
		strings.Join(m.fortuneService.GetList(), "\n"),
	)

	return []telegramclient.MessageStruct{
		telegramclient.MarkdownMessage(text),
	}, nil
}

func (m Matcher) makeRandomReplies() ([]telegramclient.MessageStruct, error) {
	fortune, err := m.fortuneService.GetRandomFortune()
	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf(
		templates.random,
		fortune.ToMarkdown(),
		fortune.File(),
	)

	return []telegramclient.MessageStruct{
		telegramclient.MarkdownMessage(text),
	}, nil
}

func (m Matcher) makeFromFileReplies(file string) ([]telegramclient.MessageStruct, error) {
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

	return []telegramclient.MessageStruct{
		telegramclient.MarkdownMessage(text),
	}, nil
}
