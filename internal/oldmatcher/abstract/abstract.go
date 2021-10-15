package abstract

import (
	"fmt"

	"github.com/br0-space/bot/internal/oldlogger"
	"github.com/br0-space/bot/internal/oldtelegram"
)

// Create a struct each matcher is inherited from
type Matcher struct{}

// Handle an error happening during matcher execution
func (m Matcher) HandleError(requestMessage oldtelegram.RequestMessage, identifier string, err error) {
	// Print error information to logger
	oldlogger.Log.Error(identifier, err.Error())

	// Send error notification to telegram
	err = oldtelegram.SendMessage(
		requestMessage,
		oldtelegram.Message{
			Text: fmt.Sprintf("⚠️ %s %s", identifier, err.Error()),
		},
	)
	if err != nil {
		oldlogger.Log.Warning(identifier, err.Error())
	}
}
