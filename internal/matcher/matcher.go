package matcher

import (
	"sync"

	"gitlab.com/br0fessional/bot/internal/logger"
	"gitlab.com/br0fessional/bot/internal/matcher/atall"
	"gitlab.com/br0fessional/bot/internal/matcher/buzzwords"
	"gitlab.com/br0fessional/bot/internal/matcher/choose"
	"gitlab.com/br0fessional/bot/internal/matcher/fortune"
	"gitlab.com/br0fessional/bot/internal/matcher/help"
	"gitlab.com/br0fessional/bot/internal/matcher/janein"
	"gitlab.com/br0fessional/bot/internal/matcher/messagestats"
	"gitlab.com/br0fessional/bot/internal/matcher/ping"
	"gitlab.com/br0fessional/bot/internal/matcher/plusplus"
	"gitlab.com/br0fessional/bot/internal/matcher/registry"
	"gitlab.com/br0fessional/bot/internal/matcher/stats"
	"gitlab.com/br0fessional/bot/internal/matcher/stonks"
	"gitlab.com/br0fessional/bot/internal/matcher/topflop"
	"gitlab.com/br0fessional/bot/internal/telegram"
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
