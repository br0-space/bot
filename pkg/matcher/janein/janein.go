package janein

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"github.com/br0-space/bot/pkg/telegram"
	"math/rand"
	"regexp"
	"strings"
)

const identifier = "janein"

var pattern = regexp.MustCompile(`(?i)^/(jn|yn)(@\w+)?($| )(.+)?$`)

var help = []interfaces.MatcherHelpStruct{{
	Command:     `jn`,
	Description: `Hilft dir, Entscheidungen zu treffen.`,
	Usage:       `/jn <Frage>`,
	Example:     `/jn ein Bier trinken`,
}}

var templates = struct {
	insult string
	yes    string
	no     string
}{
	insult: `Ob du behindert bist hab ich gefragt?\! 🤪`,
	yes:    `👍 *Ja*, du solltest *%s*\!`,
	no:     `👎 *Nein*, du solltest nicht *%s*\!`,
}

type Matcher struct {
	abstract.Matcher
}

func MakeMatcher(logger interfaces.LoggerInterface) Matcher {
	return Matcher{
		Matcher: abstract.MakeMatcher(logger, identifier, pattern, help),
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) ([]interfaces.TelegramMessageStruct, error) {
	match := m.GetCommandMatch(messageIn)
	if match == nil {
		return nil, fmt.Errorf("message does not match")
	}

	match[3] = strings.TrimSpace(match[3])
	if match[3] == "" {
		return makeReplies(templates.insult, "", messageIn.ID)
	}

	if randomYesOrNo() {
		return makeReplies(templates.yes, match[3], messageIn.ID)
	} else {
		return makeReplies(templates.no, match[3], messageIn.ID)
	}
}

func randomYesOrNo() bool {
	return rand.Float32() < 0.5
}

func makeReplies(template string, topic string, messageID int64) ([]interfaces.TelegramMessageStruct, error) {
	if strings.Contains(template, "%s") {
		template = fmt.Sprintf(
			template,
			telegram.EscapeMarkdown(topic),
		)
	}

	return []interfaces.TelegramMessageStruct{
		telegram.MakeMarkdownReply(template, messageID),
	}, nil
}
