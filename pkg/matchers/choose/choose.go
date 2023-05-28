package choose

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
)

const minOptions = 2

const identifier = "choose"

var pattern = regexp.MustCompile(`(?i)^/(choose)(@\w+)?($| )(.+)?$`)

var help = []matcher.HelpStruct{{
	Command:     `jn`,
	Description: `Wählt eine von mehreren Möglichkeiten aus.`,
	Usage:       `/choose <Option1> <Option2> <Option3>`,
	Example:     `/choose kiffen peppen nüchternBleiben`,
}}

var templates = struct {
	insult  string
	success string
}{
	insult:  `Ob du behindert bist hab ich gefragt?\! 🤪`,
	success: `👁 Das Orakel wurde befragt und hat sich entschieden für: *%s*`,
}

type Matcher struct {
	matcher.Matcher
}

func MakeMatcher() Matcher {
	return Matcher{
		Matcher: matcher.MakeMatcher(identifier, pattern, help),
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	match := m.CommandMatch(messageIn)
	if match == nil {
		return nil, fmt.Errorf("message does not match")
	}

	match[3] = strings.TrimSpace(match[3])
	if match[3] == "" {
		return makeReplies(templates.insult, "", messageIn.ID)
	}

	options := splitOptions(match[3])
	if len(options) < minOptions {
		return makeReplies(templates.insult, "", messageIn.ID)
	}

	return makeReplies(templates.success, chooseRandomOption(options), messageIn.ID)
}

func splitOptions(options string) []string {
	return strings.FieldsFunc(options, func(r rune) bool {
		return r == ' '
	})
}

func chooseRandomOption(options []string) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(options))))

	return options[int(n.Int64())]
}

func makeReplies(template string, topic string, messageID int64) ([]telegramclient.MessageStruct, error) {
	if strings.Contains(template, "%s") {
		template = fmt.Sprintf(
			template,
			telegramclient.EscapeMarkdown(topic),
		)
	}

	return []telegramclient.MessageStruct{
		telegramclient.MarkdownReply(template, messageID),
	}, nil
}
