package oldmatcher

import (
	"sync"

	"github.com/br0-space/bot/internal/oldlogger"
	"github.com/br0-space/bot/internal/oldmatcher/atall"
	"github.com/br0-space/bot/internal/oldmatcher/buzzwords"
	"github.com/br0-space/bot/internal/oldmatcher/choose"
	"github.com/br0-space/bot/internal/oldmatcher/fortune"
	"github.com/br0-space/bot/internal/oldmatcher/help"
	"github.com/br0-space/bot/internal/oldmatcher/janein"
	"github.com/br0-space/bot/internal/oldmatcher/messagestats"
	"github.com/br0-space/bot/internal/oldmatcher/music"
	"github.com/br0-space/bot/internal/oldmatcher/ping"
	"github.com/br0-space/bot/internal/oldmatcher/plusplus"
	"github.com/br0-space/bot/internal/oldmatcher/registry"
	"github.com/br0-space/bot/internal/oldmatcher/stats"
	"github.com/br0-space/bot/internal/oldmatcher/stonks"
	"github.com/br0-space/bot/internal/oldmatcher/topflop"
	"github.com/br0-space/bot/internal/oldtelegram"
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
	registry.RegisterMatcher(music.Matcher{})
	registry.RegisterMatcher(ping.Matcher{})
	registry.RegisterMatcher(plusplus.Matcher{})
	registry.RegisterMatcher(stats.Matcher{})
	registry.RegisterMatcher(stonks.Matcher{})
	registry.RegisterMatcher(topflop.Matcher{})
}

// Executes all matchers for a given request message
// Through the magic of goroutines, this is done in parallel
func ExecuteMatchers(requestMessage oldtelegram.RequestMessage) {
	oldlogger.Log.Debugf("%s wrote: %s", requestMessage.From.Username, requestMessage.Text)

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
