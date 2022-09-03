package matcher

import (
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/internal/matcher/atall"
	"github.com/br0-space/bot/internal/matcher/choose"
	"github.com/br0-space/bot/internal/matcher/janein"
	"github.com/br0-space/bot/internal/matcher/musiclinks"
	"github.com/br0-space/bot/internal/matcher/ping"
	"github.com/br0-space/bot/internal/matcher/plusplus"
	"github.com/br0-space/bot/internal/matcher/stats"
	"sync"
)

type Registry struct {
	log      interfaces.LoggerInterface
	cfg      interfaces.MatchersConfigStruct
	telegram interfaces.TelegramClientInterface
	repo     interfaces.DatabaseRepositoryInterface
	matchers []interfaces.MatcherInterface
}

func NewRegistry(
	logger interfaces.LoggerInterface,
	config interfaces.MatchersConfigStruct,
	telegram interfaces.TelegramClientInterface,
	repo interfaces.DatabaseRepositoryInterface,
) *Registry {
	registry := &Registry{
		log:      logger,
		cfg:      config,
		telegram: telegram,
		repo:     repo,
	}

	return registry
}

func (r *Registry) Init() {
	if r.cfg.Atall.Enabled {
		r.registerMatcher(atall.NewMatcher(r.log, r.cfg.Atall, r.repo.Stats()))
	}
	if r.cfg.Choose.Enabled {
		r.registerMatcher(choose.NewMatcher(r.log, r.cfg.Choose))
	}
	if r.cfg.Janein.Enabled {
		r.registerMatcher(janein.NewMatcher(r.log, r.cfg.Janein))
	}
	if r.cfg.Musiclinks.Enabled {
		r.registerMatcher(musiclinks.NewMatcher(r.log, r.cfg.Musiclinks))
	}
	if r.cfg.Ping.Enabled {
		r.registerMatcher(ping.NewMatcher(r.log, r.cfg.Ping))
	}
	if r.cfg.Plusplus.Enabled {
		r.registerMatcher(plusplus.NewMatcher(r.log, r.cfg.Plusplus, r.repo.Plusplus()))
	}
	if r.cfg.Stats.Enabled {
		r.registerMatcher(stats.NewMatcher(r.log, r.cfg.Stats, r.repo.Stats()))
	}
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

			if !m.DoesMatch(messageIn) {
				return
			}

			messagesOut, err := m.Process(messageIn)
			if err != nil {
				r.log.Error("Error while processing:", err)
				return
			}
			if messagesOut == nil {
				return
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
		messageIn.Chat.ID,
		messageIn.From.ID,
		messageIn.From.UsernameOrName(),
	); err != nil {
		r.log.Error("Error while updating user stats:", err)
	}

	if err := r.repo.MessageStats().InsertMessageStats(
		messageIn.Chat.ID,
		messageIn.From.ID,
		messageIn.WordCount(),
	); err != nil {
		r.log.Error("Error while inserting message stats:", err)
	}
}
