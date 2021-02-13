package matcher

import (
	"sync"

	"github.com/neovg/kmptnzbot/internal/logger"
	"github.com/neovg/kmptnzbot/internal/matcher/atall"
	"github.com/neovg/kmptnzbot/internal/matcher/buzzwords"
	"github.com/neovg/kmptnzbot/internal/matcher/choose"
	"github.com/neovg/kmptnzbot/internal/matcher/fortune"
	"github.com/neovg/kmptnzbot/internal/matcher/help"
	"github.com/neovg/kmptnzbot/internal/matcher/janein"
	"github.com/neovg/kmptnzbot/internal/matcher/ping"
	"github.com/neovg/kmptnzbot/internal/matcher/plusplus"
	"github.com/neovg/kmptnzbot/internal/matcher/registry"
	"github.com/neovg/kmptnzbot/internal/matcher/stats"
	"github.com/neovg/kmptnzbot/internal/matcher/stonks"
	"github.com/neovg/kmptnzbot/internal/matcher/topflop"
	"github.com/neovg/kmptnzbot/internal/telegram"
)

// At setup, create an instance of each matcher and store it in a list
func init() {
	registry.RegisterMatcher(atall.Matcher{})
	registry.RegisterMatcher(buzzwords.Matcher{})
	registry.RegisterMatcher(choose.Matcher{})
	registry.RegisterMatcher(fortune.Matcher{})
	registry.RegisterMatcher(help.Matcher{})
	registry.RegisterMatcher(janein.Matcher{})
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
