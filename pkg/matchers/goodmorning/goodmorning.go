package goodmorning

import (
	"fmt"
	"regexp"
	"time"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
)

const identifier = "goodmorning"

var pattern = regexp.MustCompile(`.+`)

var help []matcher.HelpStruct

var template = "Guten Morgen %s\\!\n\n%s\n\n_\\[from `%s`\\]_"

type Matcher struct {
	matcher.Matcher
	state   interfaces.StateServiceInterface
	fortune interfaces.FortuneServiceInterface
}

func MakeMatcher(
	state interfaces.StateServiceInterface,
	fortuneService interfaces.FortuneServiceInterface,
) Matcher {
	return Matcher{
		Matcher: matcher.MakeMatcher(identifier, pattern, help),
		state:   state,
		fortune: fortuneService,
	}
}

func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	if !m.doesMatch(messageIn) {
		return nil, nil
	}

	return m.makeReplies(messageIn)
}

func (m Matcher) doesMatch(messageIn telegramclient.WebhookMessageStruct) bool {
	now := time.Now()

	if now.Hour() < 6 || now.Hour() > 14 {
		return false
	}

	lastPost := m.state.GetLastPost(messageIn.From.ID)

	if lastPost == nil || now.Sub(*lastPost) > time.Hour*6 {
		return true
	}

	return false
}

func (m Matcher) makeReplies(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	fortune, err := m.fortune.GetRandomFortune()
	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf(
		template,
		telegramclient.EscapeMarkdown(messageIn.From.FirstnameOrUsername()),
		fortune.ToMarkdown(),
		fortune.File(),
	)

	return []telegramclient.MessageStruct{
		telegramclient.MarkdownMessage(text),
	}, nil
}
