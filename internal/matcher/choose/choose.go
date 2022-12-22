package choose

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/internal/matcher/abstract"
	"github.com/br0-space/bot/internal/telegram"
	"math/rand"
	"regexp"
	"strings"
)

const identifier = "choose"

var pattern = regexp.MustCompile(`(?i)^/(choose)(@\w+)?($| )(.+)?$`)

var help = []interfaces.MatcherHelpStruct{{
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

	match[3] = strings.TrimSpace(match[3])
	if match[3] == "" {
		return reply(templates.insult, "", messageIn.ID)
	}

	options := splitOptions(match[3])
	if len(options) < 2 {
		return reply(templates.insult, "", messageIn.ID)
	}

	return reply(templates.success, chooseOption(options), messageIn.ID)
}

func splitOptions(options string) []string {
	return strings.FieldsFunc(options, func(r rune) bool {
		return r == ' '
	})
}

func chooseOption(options []string) string {
	return options[rand.Intn(len(options))]
}

func reply(template string, topic string, messageID int64) (*[]interfaces.TelegramMessageStruct, error) {
	if strings.Contains(template, "%s") {
		template = fmt.Sprintf(
			template,
			telegram.EscapeMarkdown(topic),
		)
	}

	return &[]interfaces.TelegramMessageStruct{
		telegram.NewMarkdownReply(template, messageID),
	}, nil
}
