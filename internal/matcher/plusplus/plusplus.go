package plusplus

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/neovg/kmptnzbot/internal/db"
	"github.com/neovg/kmptnzbot/internal/matcher/abstract"
	"github.com/neovg/kmptnzbot/internal/telegram"
)

// Each matcher extends the abstract matcher
type Matcher struct{
	abstract.Matcher
}

type Token struct {
	name string
	increment int
}

// Return the identifier of this matcher for use in logging
func (m Matcher) Identifier() string {
	return "plusplus"
}

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
	// Tokenize text
	matches := m.getMatches(requestMessage.TextOrCaption())
	tokens := m.getTokens(matches)

	// Process tokens
	m.processTokens(requestMessage, tokens)

	return nil
}

// Return a list of substrings in the format {name}++ / {name}+- / {name}-- contained in a text
// The list is not unique, the same substring might be contained multiple times
func (m Matcher) getMatches(text string) []string {
	// Check if message starts with / and if yes, ignore it
	match, _ := regexp.MatchString(`^/`, text)
	if match {
		return make([]string, 0, 0)
	}

	// Initialize the regular expression
	r := regexp.MustCompile(`(^|\s)[\p{L}\w]+(\+\+|\+-|--)`)

	// Find all occurrences of {name}++ / {name}+- / {name}--
	matches := r.FindAllString(text, -1)

	// Trim matches to get rid of leading spaces
	for i := range matches {
		matches[i] = strings.TrimSpace(matches[i])
	}

	return matches
}

// Take a list of substrings in the format {name}++ / {name}+- / {name}-- and return a map of tokens
// Each token consists of a name and a value determined by the number matches with this name and ++, +- or --
func (m Matcher) getTokens(matches []string) []Token {
	// Initialize a map for more performant collection of values
	values := make(map[string]int)

	// Add each token to the map with its name as key and set its value to the sum of ++, +- and -- for this name
	for _, match := range matches {
		// Ignore the case to avoid duplicates
		name := strings.ToLower(match[:len(match)-2])

		// ++ will increment by 1, -- will decrement and +- leaves the value as it is
		increment := 0
		switch match[len(match)-2:] {
		case "++":
			increment = 1
		case "--":
			increment = -1
		}

		// If the token was not yet added to the map, add it with the initial value
		// Otherwise update the value of the existing entry
		if _, exists := values[name]; !exists {
			values[name] = increment
		} else {
			values[name] = values[name] + increment
		}
	}

	// Transform the map into a slice of token objects
	tokens := make([]Token, 0, len(values))
	for name, increment := range values {
		tokens = append(tokens, Token{name: name, increment: increment})
	}

	return tokens
}

// Take a list of tokens and process each one in a goroutine
func (m Matcher) processTokens(requestMessage telegram.RequestMessage, tokens []Token) {
	// Create a wait group for synchronization
	var waitGroup sync.WaitGroup

	// We need to wait until all tokens are processed
	waitGroup.Add(len(tokens))

	// Launch a goroutine for each token
	for _, token := range tokens {
		go func(token Token) {
			defer waitGroup.Done()

			// Process the token
			err := m.processToken(requestMessage, token)
			if err != nil {
				m.HandleError(requestMessage, m.Identifier(), err)
			}
		}(token)
	}

	// Wait until all tokens are processed
	waitGroup.Wait()
}

// Update a token in the database and send a notification message with the new value to Telegram
func (m Matcher) processToken(requestMessage telegram.RequestMessage, token Token) error {
	// Update the database
	newValue := db.IncrementPlusplus(token.name, token.increment)

	// Send message to Telegram
	return m.sendResponse(requestMessage, token, newValue)
}

// Send a message with the new value to Telegram
func (m Matcher) sendResponse(requestMessage telegram.RequestMessage, token Token, newValue int) error {
	mode := "+-"
	if newValue > 0 {
		mode = "++"
	}
	if newValue < 0 {
		mode = "--"
	}

	responseMessage := telegram.Message{
		Text: fmt.Sprintf(
			"\\[%s] *%s* ist jetzt auf *%d*",
			mode,
			token.name,
			newValue,
		),
		ParseMode: "Markdown",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}
