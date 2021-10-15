package help

import (
	"fmt"
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
	return "help"
}

// This is a command matcher and generates a help item
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{{
		Command:     "help",
		Description: "Zeigt die verfügbaren Befehle an",
	}}
}

// Process a message received from telegram
func (m Matcher) ProcessRequestMessage(requestMessage oldtelegram.RequestMessage) error {
	// Check if text starts with /help and if not, ignore it
	if doesMatch := m.doesMatch(requestMessage.Text); !doesMatch {
		return nil
	}

	// Choose one option and send the result
	return m.sendResponse(requestMessage)
}

// Check if a text starts with /help
func (m Matcher) doesMatch(text string) bool {
	// Check if message starts with /help
	match, _ := regexp.MatchString(`^/help(@|\s|$)`, text)

	return match
}

// Send the result to the user who sent the request message
func (m Matcher) sendResponse(requestMessage oldtelegram.RequestMessage) error {
	text := ""

	for _, matcher := range registry.GetRegisteredMatchers() {
		for _, helpItem := range matcher.GetHelpItems() {
			text = text + fmt.Sprintf("%s - %s\n", helpItem.Command, helpItem.Description)
		}
	}

	responseMessage := oldtelegram.Message{
		Text: text,
	}

	return oldtelegram.SendMessage(requestMessage, responseMessage)
}
