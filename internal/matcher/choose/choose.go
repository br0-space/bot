package choose

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/neovg/kmptnzbot/internal/matcher/abstract"
	"github.com/neovg/kmptnzbot/internal/matcher/registry"
	"github.com/neovg/kmptnzbot/internal/telegram"
)

// Each matcher extends the abstract matcher
type Matcher struct {
	abstract.Matcher
}

// Return the identifier of this matcher for use in logging
func (m Matcher) Identifier() string {
	return "choose"
}

// This is a command matcher and generates a help item
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{{
		Command:     "choose",
		Description: "Wählt eine aus mehreren Optionen (Beispiel: `/choose one two three`)",
	}}
}

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
	// Check if text starts with /choose and if not, ignore it
	if doesMatch := m.doesMatch(requestMessage.Text); !doesMatch {
		return nil
	}

	// Get options to choose from
	options := m.getOptions(requestMessage.Text)

	// If not enough options were found, insult the idiot who sent the request message
	if len(options) < 2 {
		return m.sendInsultResponse(requestMessage)
	}

	// Choose one option and send the result
	return m.sendResultResponse(requestMessage, m.getRandomOption(options))
}

// Check if a text starts with /choose
func (m Matcher) doesMatch(text string) bool {
	// Check if message starts with /choose
	match, _ := regexp.MatchString(`^/choose(@|\s|$)`, text)

	return match
}

// Returns a slice of all options to choose from
// An option is any word prepended by a whitespace (to ignore the command itself)
func (m Matcher) getOptions(text string) []string {
	// Initialize the regular expression
	r := regexp.MustCompile(`\s\S+`)

	// Find all words to choose from
	words := r.FindAllString(text, -1)

	// Trim words to get rid of leading spaces
	for i := range words {
		words[i] = strings.TrimSpace(words[i])
	}

	return words
}

// Returns a random element from a slice of strings
func (m Matcher) getRandomOption(options []string) string {
	return options[rand.Intn(len(options))]
}

// Send an insult to the user who sent the request message
func (m Matcher) sendInsultResponse(requestMessage telegram.RequestMessage) error {
	responseMessage := telegram.Message{
		Text:             "Ob du behindert bist hab ich gefragt?! 🤪",
		ReplyToMessageID: requestMessage.ID,
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}

// Send the result to the user who sent the request message
func (m Matcher) sendResultResponse(requestMessage telegram.RequestMessage, result string) error {
	responseMessage := telegram.Message{
		Text:             fmt.Sprintf("👁 Das Orakel wurde befragt und hat sich entschieden für: %s", result),
		ReplyToMessageID: requestMessage.ID,
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}
