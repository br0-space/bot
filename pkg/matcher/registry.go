package matcher

import (
	"fmt"
	logger "github.com/br0-space/bot-logger"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/atall"
	"github.com/br0-space/bot/pkg/matcher/buzzwords"
	"github.com/br0-space/bot/pkg/matcher/choose"
	"github.com/br0-space/bot/pkg/matcher/fortune"
	"github.com/br0-space/bot/pkg/matcher/goodmorning"
	"github.com/br0-space/bot/pkg/matcher/janein"
	"github.com/br0-space/bot/pkg/matcher/musiclinks"
	"github.com/br0-space/bot/pkg/matcher/ping"
	"github.com/br0-space/bot/pkg/matcher/plusplus"
	"github.com/br0-space/bot/pkg/matcher/stats"
	"github.com/br0-space/bot/pkg/matcher/topflop"
	"github.com/br0-space/bot/pkg/matcher/xkcd"
	"github.com/br0-space/bot/pkg/telegram"
	"sync"
)

const errorTemplate = "⚠️ *Error in matcher \"%s\"*\n\n%s"

type Registry struct {
	log              logger.Interface
	state            interfaces.StateServiceInterface
	telegram         interfaces.TelegramClientInterface
	messageStatsRepo interfaces.MessageStatsRepoInterface
	plusplusRepo     interfaces.PlusplusRepoInterface
	userStatsRepo    interfaces.UserStatsRepoInterface
	fortuneService   interfaces.FortuneServiceInterface
	songlinkService  interfaces.SonglinkServiceInterface
	xkcdService      interfaces.XkcdServiceInterface
	matchers         []interfaces.MatcherInterface
}

func NewRegistry(
	state interfaces.StateServiceInterface,
	telegram interfaces.TelegramClientInterface,
	messageStatsRepo interfaces.MessageStatsRepoInterface,
	plusplusRepo interfaces.PlusplusRepoInterface,
	userStatsRepo interfaces.UserStatsRepoInterface,
	fortuneService interfaces.FortuneServiceInterface,
	songlinkService interfaces.SonglinkServiceInterface,
	xkcdService interfaces.XkcdServiceInterface,
) *Registry {
	registry := &Registry{
		log:              logger.New(),
		state:            state,
		telegram:         telegram,
		messageStatsRepo: messageStatsRepo,
		plusplusRepo:     plusplusRepo,
		userStatsRepo:    userStatsRepo,
		fortuneService:   fortuneService,
		songlinkService:  songlinkService,
		xkcdService:      xkcdService,
	}

	return registry
}

func (r *Registry) Init() {
	r.registerMatcher(atall.MakeMatcher(r.userStatsRepo))
	r.registerMatcher(buzzwords.MakeMatcher(r.plusplusRepo))
	r.registerMatcher(choose.MakeMatcher())
	r.registerMatcher(goodmorning.MakeMatcher(r.state, r.fortuneService))
	r.registerMatcher(fortune.MakeMatcher(r.fortuneService))
	r.registerMatcher(janein.MakeMatcher())
	r.registerMatcher(musiclinks.MakeMatcher(r.songlinkService))
	r.registerMatcher(ping.MakeMatcher())
	r.registerMatcher(plusplus.MakeMatcher(r.plusplusRepo))
	r.registerMatcher(stats.MakeMatcher(r.userStatsRepo))
	r.registerMatcher(topflop.MakeMatcher(r.plusplusRepo))
	r.registerMatcher(xkcd.MakeMatcher(r.xkcdService))
}

func (r *Registry) registerMatcher(matcher interfaces.MatcherInterface) {
	r.log.Debug("Registering matcher", matcher.Identifier())

	r.matchers = append(r.matchers, matcher)
}

func (r *Registry) Process(messageIn interfaces.TelegramWebhookMessageStruct) {
	r.log.Debugf("Processing message from %s: %s", messageIn.From.Username, messageIn.Text)

	// Create a wait group for synchronization
	var waitGroup sync.WaitGroup

	// We need to wait until all matchers are executed
	waitGroup.Add(len(r.matchers))

	// Launch a goroutine for each matcher
	for _, m := range r.matchers {
		go func(m interfaces.MatcherInterface) {
			defer waitGroup.Done()

			if !m.IsEnabled() {
				return
			}

			if !m.DoesMatch(messageIn) {
				return
			}

			messagesOut, err := m.Process(messageIn)
			if messagesOut == nil {
				messagesOut = []interfaces.TelegramMessageStruct{}
			}
			if err != nil {
				r.log.Errorf("Error in matcher %s: %s", m.Identifier(), err)

				messagesOut = append(
					messagesOut,
					telegram.MakeMarkdownReply(
						fmt.Sprintf(
							errorTemplate,
							m.Identifier(),
							telegram.EscapeMarkdown(err.Error()),
						),
						messageIn.ID,
					),
				)
			}

			for _, messageOut := range messagesOut {
				if err := r.telegram.SendMessage(messageIn.Chat.ID, messageOut); err != nil {
					r.log.Error("Error while sending message:", err)
				}
			}
		}(m)
	}

	waitGroup.Wait()

}
