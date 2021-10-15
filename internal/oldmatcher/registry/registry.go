package registry

import (
	"github.com/br0-space/bot/internal/oldlogger"
	"github.com/br0-space/bot/internal/oldtelegram"
)

// Each matcher must implement a function to process request messages
type Matcher interface {
	Identifier() string
	GetHelpItems() []HelpItem
	ProcessRequestMessage(requestMessage oldtelegram.RequestMessage) error
	HandleError(requestMessage oldtelegram.RequestMessage, identifier string, err error)
}

type HelpItem struct {
	Command     string
	Description string
}

// List of all registered matcher instances
var matchers = make([]Matcher, 0)

// Add a matcher to the list
func RegisterMatcher(matcher Matcher) {
	oldlogger.Log.Debug("register matcher", matcher.Identifier())

	matchers = append(matchers, matcher)
}

// Return all registered matchers
func GetRegisteredMatchers() []Matcher {
	return matchers
}
