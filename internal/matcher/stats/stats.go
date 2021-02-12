package stats

import (
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
	return "choose"
}

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
	db.UpdateStats(requestMessage.From.ID, requestMessage.From.Username)

	return nil
}