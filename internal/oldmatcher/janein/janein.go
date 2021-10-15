package janein

import (
	"fmt"
	"math/rand"
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
	return "janein"
}

// This is a command matcher and generates a help item
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{{
		Command:     "jn",
		Description: "Sagt dir, ob du etwas machen sollst oder nicht (Beispiel: `/jn Bugs fixen`)",
	}}
}

// Process a message received from telegram
func (m Matcher) ProcessRequestMessage(requestMessage oldtelegram.RequestMessage) error {
	// Check if text starts with /jn or /yn and if not, ignore it
	if doesMatch := m.doesMatch(requestMessage.Text); !doesMatch {
		return nil
	}

	// Get the option
	option := m.getOption(requestMessage.Text)

	// If not enough options were found, insult the idiot who sent the request message
	if len(option) == 0 {
		return m.sendInsultResponse(requestMessage)
	}

	// Choose one option and send the result
	return m.sendResultResponse(requestMessage, option, m.getRandomYesOrNo())
}

// Check if a text starts with /jn or /yn
func (m Matcher) doesMatch(text string) bool {
	match, _ := regexp.MatchString(`^/(jn|yn)(@|\s|$)`, text)

	return match
}

// Check if a text starts with /jn or /yn and return the text behind
func (m Matcher) getOption(text string) string {
	match, _ := regexp.MatchString(`^/(jn|yn) .+`, text)
	if !match {
		return ""
	}

	return text[4:]
}

// Maybe yeeeees, maybe noooooo
func (m Matcher) getRandomYesOrNo() bool {
	return rand.Float32() < 0.5
}

// Send an insult to the user who sent the request message
func (m Matcher) sendInsultResponse(requestMessage oldtelegram.RequestMessage) error {
	responseMessage := oldtelegram.Message{
		Text:             "Ob du behindert bist hab ich gefragt?! 🤪",
		ReplyToMessageID: requestMessage.ID,
	}

	return oldtelegram.SendMessage(requestMessage, responseMessage)
}

// Send a message with the result to telegram
func (m Matcher) sendResultResponse(requestMessage oldtelegram.RequestMessage, text string, result bool) error {
	if result {
		text = fmt.Sprintf("👍 Ja, du solltest %s!", text)
	} else {
		text = fmt.Sprintf("👎 Nein, du solltest nicht %s!", text)
	}

	responseMessage := oldtelegram.Message{
		Text:             text,
		ReplyToMessageID: requestMessage.ID,
	}

	return oldtelegram.SendMessage(requestMessage, responseMessage)
}
