package matcher

import (
	"sync"

	"github.com/kmptnz/bot/internal/logger"
	"github.com/kmptnz/bot/internal/matcher/atall"
	"github.com/kmptnz/bot/internal/matcher/buzzwords"
	"github.com/kmptnz/bot/internal/matcher/choose"
	"github.com/kmptnz/bot/internal/matcher/fortune"
	"github.com/kmptnz/bot/internal/matcher/help"
	"github.com/kmptnz/bot/internal/matcher/janein"
	"github.com/kmptnz/bot/internal/matcher/messagestats"
	"github.com/kmptnz/bot/internal/matcher/ping"
	"github.com/kmptnz/bot/internal/matcher/plusplus"
	"github.com/kmptnz/bot/internal/matcher/registry"
	"github.com/kmptnz/bot/internal/matcher/stats"
	"github.com/kmptnz/bot/internal/matcher/stonks"
	"github.com/kmptnz/bot/internal/matcher/topflop"
	"github.com/kmptnz/bot/internal/telegram"
)

// At setup, create an instance of each matcher and store it in a list
func init() {
	registry.RegisterMatcher(atall.Matcher{})
	registry.RegisterMatcher(buzzwords.Matcher{})
	registry.RegisterMatcher(choose.Matcher{})
	registry.RegisterMatcher(fortune.Matcher{})
	registry.RegisterMatcher(help.Matcher{})
	registry.RegisterMatcher(janein.Matcher{})
	registry.RegisterMatcher(messagestats.Matcher{})
	registry.RegisterMatcher(ping.Matcher{})
	registry.RegisterMatcher(plusplus.Matcher{})
	registry.RegisterMatcher(stats.Matcher{})
	registry.RegisterMatcher(stonks.Matcher{})
	registry.RegisterMatcher(topflop.Matcher{})
}

// Executes all matchers for a given request message
// Through the magic of goroutines, this is done in parallel
func ExecuteMatchers(requestMessage telegram.RequestMessage) {
	logger.Log.Debugf("%s wrote: %s", requestMessage.From.Username, requestMessage.Text)

	// Create a wait group for synchronization
	var waitGroup sync.WaitGroup

	// We need to wait until all matchers are executed
	waitGroup.Add(len(registry.GetRegisteredMatchers()))

	// Launch a goroutine for each matcher
	for _, matcher := range registry.GetRegisteredMatchers() {
		go func(matcher registry.Matcher) {
			defer waitGroup.Done()

			// Let the matcher process the request message
			if err := matcher.ProcessRequestMessage(requestMessage); err != nil {
				matcher.HandleError(requestMessage, matcher.Identifier(), err)
				return
			}
		}(matcher)
	}

	// Wait until all matchers are executed
	waitGroup.Wait()
}
