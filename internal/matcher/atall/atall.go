package atall

import (
	"regexp"
	"strings"

	"github.com/neovg/kmptnzbot/internal/db"
	"github.com/neovg/kmptnzbot/internal/matcher/abstract"
	"github.com/neovg/kmptnzbot/internal/telegram"
)

// Each matcher extends the abstract matcher
type Matcher struct {
	abstract.Matcher
}

// Return the identifier of this matcher for use in logging
func (m Matcher) Identifier() string {
	return "atall"
}

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
	// Check if text starts with /ping and if not, ignore it
	if doesMatch := m.doesMatch(requestMessage.Text); !doesMatch {
		return nil
	}

	// Choose one option and send the result
	return m.sendResponse(requestMessage)
}

// Check if a text starts with /ping
func (m Matcher) doesMatch(text string) bool {
	// Check if message is a command and if yes, ignore ir
	cmd, _ := regexp.MatchString(`^/`, text)
	if cmd {
		return false
	}

	// Check if message contains @all or @alle
	match, _ := regexp.MatchString(`(^|\s)@alle?(\s|$)`, text)

	return match
}

// Send the result to the user who sent the request message
func (m Matcher) sendResponse(requestMessage telegram.RequestMessage) error {
	usernames := db.FindAllUsernames(requestMessage.From.Username)

	responseMessage := telegram.Message{
		Text:             strings.Join(usernames, " "),
		ReplyToMessageID: requestMessage.ID,
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}