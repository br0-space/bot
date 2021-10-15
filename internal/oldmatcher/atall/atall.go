package atall

import (
	"regexp"
	"strings"

	"github.com/br0-space/bot/internal/db"
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
	return "atall"
}

// This matcher is no command and generates no help items
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{}
}

// Process a message received from telegram
func (m Matcher) ProcessRequestMessage(requestMessage oldtelegram.RequestMessage) error {
	// Check if text contains @all or @alle and if not, ignore it
	if doesMatch := m.doesMatch(requestMessage.Text); !doesMatch {
		return nil
	}

	// Choose one option and send the result
	return m.sendResponse(requestMessage)
}

// Check if a text contains @all or @alle
func (m Matcher) doesMatch(text string) bool {
	// Check if message is a command and if yes, ignore it
	cmd, _ := regexp.MatchString(`^/`, text)
	if cmd {
		return false
	}

	// Check if message contains @all or @alle
	match, _ := regexp.MatchString(`(^|\s)@alle?(\s|$)`, text)

	return match
}

// Send the original text together with a list of mentioned users
func (m Matcher) sendResponse(requestMessage oldtelegram.RequestMessage) error {
	usernames := db.FindAllUsernames(requestMessage.From.Username)

	text := requestMessage.TextOrCaption()
	text = strings.ReplaceAll(text, "@alle", "")
	text = strings.ReplaceAll(text, "@all", "")
	text = text + " " + strings.Join(usernames, " ")

	responseMessage := oldtelegram.Message{
		Text:             text,
		ReplyToMessageID: requestMessage.ID,
	}

	return oldtelegram.SendMessage(requestMessage, responseMessage)
}
