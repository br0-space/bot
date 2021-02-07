package matcher

import (
	"sync"

	"github.com/neovg/kmptnzbot/internal/logger"
	"github.com/neovg/kmptnzbot/internal/matcher/choose"
	"github.com/neovg/kmptnzbot/internal/matcher/janein"
	"github.com/neovg/kmptnzbot/internal/matcher/ping"
	"github.com/neovg/kmptnzbot/internal/matcher/plusplus"
	"github.com/neovg/kmptnzbot/internal/matcher/stonks"
	"github.com/neovg/kmptnzbot/internal/telegram"
)

// Each matcher must implement a function to process request messages
type Matcher interface {
	ProcessRequestMessage(requestMessage telegram.RequestMessage) error
	HandleError(requestMessage telegram.RequestMessage, identifier string, err error)
	Identifier() string
}

// List of all registered matcher instances
var matchers = make([]Matcher, 0, 0)

// At setup, create an instance of each matcher and store it in a list
func init() {
	RegisterMatchers()
}

// Adds a matcher to the list
func registerMatcher(matcher Matcher) {
	logger.Log.Debug("register matcher", matcher.Identifier())

	matchers = append(matchers, matcher)
}

// Creates an instance of each matcher and adds it to the list
func RegisterMatchers() {
	registerMatcher(choose.Matcher{})
	registerMatcher(janein.Matcher{})
	registerMatcher(ping.Matcher{})
	registerMatcher(plusplus.Matcher{})
	registerMatcher(stonks.Matcher{})
}

// Executes all matchers for a given request message
// Through the magic of goroutines, this is done in parallel
func ExecuteMatchers(requestMessage telegram.RequestMessage) {
	logger.Log.Infof("%s wrote: %s", requestMessage.From.Username, requestMessage.Text)

	// Create a wait group for synchronization
	var waitGroup sync.WaitGroup

	// We need to wait until all matchers are executed
	waitGroup.Add(len(matchers))

	// Launch a goroutine for each matcher
	for _, matcher := range matchers {
		go func(matcher Matcher) {
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

	logger.Log.Debug("all matchers executed")
}
