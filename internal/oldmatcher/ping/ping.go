package ping

import (
	"regexp"

	"github.com/br0-space/bot/internal/oldmatcher/abstract"
	"github.com/br0-space/bot/internal/oldmatcher/registry"
	"github.com/br0-space/bot/internal/oldtelegram"
)

// Each matcher extends the abstract matcher
type Matcher struct {
	abstract.Matcher
}

// Return the identifier of this matcher for use in logging
func (m Matcher) Identifier() string {
	return "ping"
}

// This is a command matcher and generates a help item
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{{
		Command:     "ping",
		Description: "Antwortet mit `pong`",
	}}
}

// Process a message received from telegram
func (m Matcher) ProcessRequestMessage(requestMessage oldtelegram.RequestMessage) error {
	// Check if text starts with /ping and if not, ignore it
	if doesMatch := m.doesMatch(requestMessage.Text); !doesMatch {
		return nil
	}

	// Choose one option and send the result
	return m.sendResponse(requestMessage)
}

// Check if a text starts with /ping
func (m Matcher) doesMatch(text string) bool {
	// Check if message starts with /ping
	match, _ := regexp.MatchString(`^/ping(@|\s|$)`, text)

	return match
}

// Send the result to the user who sent the request message
func (m Matcher) sendResponse(requestMessage oldtelegram.RequestMessage) error {
	responseMessage := oldtelegram.Message{
		Text:             "pong",
		ReplyToMessageID: requestMessage.ID,
	}

	return oldtelegram.SendMessage(requestMessage, responseMessage)
}
