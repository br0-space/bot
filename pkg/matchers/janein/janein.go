package janein

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
)

const randMax = 1000

const identifier = "janein"

var pattern = regexp.MustCompile(`(?i)^/(jn|yn)(@\w+)?($| )(.+)?$`)

var help = []matcher.HelpStruct{{
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

	if randomYesOrNo() {
		return makeReplies(templates.yes, match[3], messageIn.ID)
	} else {
		return makeReplies(templates.no, match[3], messageIn.ID)
	}
}

func randomYesOrNo() bool {
	n, _ := rand.Int(rand.Reader, big.NewInt(randMax))

	return n.Int64() < randMax/2
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
