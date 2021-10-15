package registry

import (
	"sync"

	"github.com/br0-space/bot/container"
	"github.com/br0-space/bot/internal/logger"
	"github.com/br0-space/bot/internal/matcher"
	"github.com/br0-space/bot/internal/matcher/ping"
	"github.com/br0-space/bot/internal/telegram/webhook"
)

type Registry struct {
	log logger.Interface
	matchers []matcher.Interface
}

var registry *Registry

func NewRegistry() *Registry {
	if registry == nil {
		registry = &Registry{
			log:      container.ProvideLoggerService(),
			matchers: make([]matcher.Interface, 0),
		}
		registry.registerMatcher(ping.MakeMatcher())
	}

	return registry
}

func (r *Registry) registerMatcher(matcher matcher.Interface) {
	r.matchers = append(r.matchers, matcher)
}

// Executes all matchers for a given request message
// Through the magic of goroutines, this is done in parallel
func (r *Registry) ProcessWebhookMessageInAllMatchers(messageIn webhook.Message) {
	r.log.Debugf("%s wrote: %s", messageIn.From.Username, messageIn.Text)

	// Create a wait group for synchronization
	var waitGroup sync.WaitGroup

	// We need to wait until all matchers are executed
	waitGroup.Add(len(r.matchers))

	// Launch a goroutine for each matcher
	for _, m := range r.matchers {
		go func(m matcher.Interface) {
			defer waitGroup.Done()

			// Let the matcher process the request message
			if err := m.ProcessMessage(messageIn); err != nil {
				m.HandleError(messageIn, err)
				return
			}
		}(m)
	}
}