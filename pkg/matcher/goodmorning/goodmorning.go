package goodmorning

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/abstract"
	"github.com/br0-space/bot/pkg/telegram"
	"regexp"
	"time"
)

const identifier = "goodmorning"

var pattern = regexp.MustCompile(`.+`)

var help []interfaces.MatcherHelpStruct

var template = "Guten Morgen %s\\!\n\n%s\n\n_\\[from `%s`\\]_"

type Matcher struct {
	abstract.Matcher
	state   interfaces.StateServiceInterface
	fortune interfaces.FortuneServiceInterface
}

func MakeMatcher(
	logger interfaces.LoggerInterface,
	state interfaces.StateServiceInterface,
	fortuneService interfaces.FortuneServiceInterface,
) Matcher {
	return Matcher{
		Matcher: abstract.MakeMatcher(logger, identifier, pattern, help),
		state:   state,
		fortune: fortuneService,
	}
}

func (m Matcher) Process(messageIn interfaces.TelegramWebhookMessageStruct) ([]interfaces.TelegramMessageStruct, error) {
	if !m.doesMatch(messageIn) {
		return nil, nil
	}

	return m.makeReplies(messageIn)
}

func (m Matcher) doesMatch(messageIn interfaces.TelegramWebhookMessageStruct) bool {
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

func (m Matcher) makeReplies(messageIn interfaces.TelegramWebhookMessageStruct) ([]interfaces.TelegramMessageStruct, error) {
	fortune, err := m.fortune.GetRandomFortune()
	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf(
		template,
		telegram.EscapeMarkdown(messageIn.From.FirstnameOrUsername()),
		fortune.ToMarkdown(),
		fortune.File(),
	)

	return []interfaces.TelegramMessageStruct{
		telegram.MakeMarkdownMessage(text),
	}, nil
}
