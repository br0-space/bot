package janein

import (
	"fmt"
	"math/rand"
	"regexp"

	"github.com/neovg/kmptnzbot/internal/matcher/abstract"
	"github.com/neovg/kmptnzbot/internal/telegram"
)

// Each matcher extends the abstract matcher
type Matcher struct {
	abstract.Matcher
}

// Return the identifier of this matcher for use in logging
func (m Matcher) Identifier() string {
	return "janein"
}

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
	// Check if text starts with /jn
	match := m.getMatch(requestMessage.Text)

	// If match is empty, text didn't start with /jn
	if match == "" {
		return nil
	}

	// Send a randomized response
	return m.sendResponse(requestMessage, match, m.getRandomYesOrNo())
}

// Check if a text starts with /jn and return the text behind
func (m Matcher) getMatch(text string) string {
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

// Send a message with the result to Telegram
func (m Matcher) sendResponse(requestMessage telegram.RequestMessage, text string, result bool) error {
	if result {
		text = fmt.Sprintf("👍 Ja, du solltest %s!", text)
	} else {
		text = fmt.Sprintf("👎 Nein, du solltest nicht %s!", text)
	}

	responseMessage := telegram.Message{
		Text:             text,
		ReplyToMessageID: requestMessage.ID,
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}
