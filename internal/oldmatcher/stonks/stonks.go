package stonks

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

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
	return "stonks"
}

// This matcher is no command and generates no help items
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{}
}

// Process a message received from telegram
func (m Matcher) ProcessRequestMessage(requestMessage oldtelegram.RequestMessage) error {
	// Extract symbols from text
	symbols := m.getSymbols(requestMessage.TextOrCaption())
	if len(symbols) == 0 {
		return nil
	}

	// Load quotes for symbols
	quotes, err := getQuotes(symbols)
	if err != nil {
		m.HandleError(requestMessage, m.Identifier(), err)
	}

	// Process quotes
	m.processQuotes(requestMessage, quotes)

	return nil
}

// Return a list of stonk symbols contained in a text
func (m Matcher) getSymbols(text string) []string {
	// Check if message starts with / and if yes, ignore it
	match, _ := regexp.MatchString(`^/`, text)
	if match {
		return make([]string, 0)
	}

	// Initialize the regular expression
	r := regexp.MustCompile(`(^|\s)\$[A-Z0-9:.]+`)

	// Find all occurrences of ${symbol}
	symbols := r.FindAllString(text, -1)

	// Trim matches to get rid of leading spaces and the dollar sign
	for i := range symbols {
		symbols[i] = strings.TrimSpace(symbols[i])
		symbols[i] = strings.TrimLeft(symbols[i], "$")
	}

	return symbols
}

// Take a list of stonk quotes and process each one in a goroutine
func (m Matcher) processQuotes(requestMessage oldtelegram.RequestMessage, quotes []Quote) {
	// Create a wait group for synchronization
	var waitGroup sync.WaitGroup

	// We need to wait until all quotes are processed
	waitGroup.Add(len(quotes))

	// Launch a goroutine for each quote
	for _, quote := range quotes {
		go func(quote Quote) {
			defer waitGroup.Done()

			// Process the token
			err := m.sendResponse(requestMessage, quote)
			if err != nil {
				m.HandleError(requestMessage, m.Identifier(), err)
			}
		}(quote)
	}

	// Wait until all quotes are processed
	waitGroup.Wait()
}

// Send a message with the stonk quote to telegram
func (m Matcher) sendResponse(requestMessage oldtelegram.RequestMessage, quote Quote) error {
	changeEmoji := "🤷"
	if quote.Change > 0 {
		changeEmoji = "🚀"
	}
	if quote.Change < 0 {
		changeEmoji = "🔥"
	}

	responseMessage := oldtelegram.Message{
		Text: fmt.Sprintf(
			"%s %s (%s): %.2f %s (%.2f%%)",
			changeEmoji,
			quote.ResponseQuote.ShortName,
			quote.ResponseQuote.ExchangeName,
			quote.Price,
			quote.ResponseQuote.Currency,
			quote.ChangePercent,
		),
		ParseMode: "Markdown",
	}

	return oldtelegram.SendMessage(requestMessage, responseMessage)
}
