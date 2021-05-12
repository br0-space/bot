package matcher

import (
	"sync"

	"github.com/br0-space/bot/internal/logger"
	"github.com/br0-space/bot/internal/matcher/atall"
	"github.com/br0-space/bot/internal/matcher/buzzwords"
	"github.com/br0-space/bot/internal/matcher/choose"
	"github.com/br0-space/bot/internal/matcher/fortune"
	"github.com/br0-space/bot/internal/matcher/help"
	"github.com/br0-space/bot/internal/matcher/janein"
	"github.com/br0-space/bot/internal/matcher/messagestats"
	"github.com/br0-space/bot/internal/matcher/ping"
	"github.com/br0-space/bot/internal/matcher/plusplus"
	"github.com/br0-space/bot/internal/matcher/registry"
	"github.com/br0-space/bot/internal/matcher/stats"
	"github.com/br0-space/bot/internal/matcher/stonks"
	"github.com/br0-space/bot/internal/matcher/topflop"
	"github.com/br0-space/bot/internal/telegram"
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
