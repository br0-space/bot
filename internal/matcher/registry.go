package matcher

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/internal/matcher/atall"
	"github.com/br0-space/bot/internal/matcher/buzzwords"
	"github.com/br0-space/bot/internal/matcher/choose"
	"github.com/br0-space/bot/internal/matcher/fortune"
	"github.com/br0-space/bot/internal/matcher/janein"
	"github.com/br0-space/bot/internal/matcher/musiclinks"
	"github.com/br0-space/bot/internal/matcher/ping"
	"github.com/br0-space/bot/internal/matcher/plusplus"
	"github.com/br0-space/bot/internal/matcher/stats"
	"github.com/br0-space/bot/internal/matcher/topflop"
	"github.com/br0-space/bot/internal/telegram"
	"sync"
)

const errorTemplate = "⚠️ *Error in matcher \"%s\"*\n\n%s"

type Registry struct {
	log      interfaces.LoggerInterface
	telegram interfaces.TelegramClientInterface
	repo     interfaces.DatabaseRepositoryInterface
	matchers []interfaces.MatcherInterface
}

func NewRegistry(
	logger interfaces.LoggerInterface,
	telegram interfaces.TelegramClientInterface,
	repo interfaces.DatabaseRepositoryInterface,
) *Registry {
	registry := &Registry{
		log:      logger,
		telegram: telegram,
		repo:     repo,
	}

	return registry
}

func (r *Registry) Init() {
	r.registerMatcher(atall.NewMatcher(r.log, r.repo.Stats()))
	r.registerMatcher(buzzwords.NewMatcher(r.log, r.repo.Plusplus()))
	r.registerMatcher(choose.NewMatcher(r.log))
	r.registerMatcher(fortune.NewMatcher(r.log))
	r.registerMatcher(janein.NewMatcher(r.log))
	r.registerMatcher(musiclinks.NewMatcher(r.log))
	r.registerMatcher(ping.NewMatcher(r.log))
	r.registerMatcher(plusplus.NewMatcher(r.log, r.repo.Plusplus()))
	r.registerMatcher(stats.NewMatcher(r.log, r.repo.Stats()))
	r.registerMatcher(topflop.NewMatcher(r.log, r.repo.Plusplus()))
}

func (r *Registry) registerMatcher(matcher interfaces.MatcherInterface) {
	r.log.Debug("Registering matcher", matcher.GetIdentifier())

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
				messagesOut = &[]interfaces.TelegramMessageStruct{}
			}
			if err != nil {
				r.log.Errorf("Error in matcher %s: %s", m.GetIdentifier(), err)

				*messagesOut = append(
					*messagesOut,
					telegram.NewMarkdownReply(
						fmt.Sprintf(
							errorTemplate,
							m.GetIdentifier(),
							telegram.EscapeMarkdown(err.Error()),
						),
						messageIn.ID,
					),
				)
			}

			for _, messageOut := range *messagesOut {
				if err := r.telegram.SendMessage(messageIn.Chat.ID, messageOut); err != nil {
					r.log.Error("Error while sending message:", err)
				}
			}
		}(m)
	}

	waitGroup.Wait()

	if err := r.repo.Stats().UpdateStats(
		messageIn.From.ID,
		messageIn.From.UsernameOrName(),
	); err != nil {
		r.log.Error("Error while updating user stats:", err)
	}

	if err := r.repo.MessageStats().InsertMessageStats(
		messageIn.From.ID,
		messageIn.WordCount(),
	); err != nil {
		r.log.Error("Error while inserting message stats:", err)
	}
}
