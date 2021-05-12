package messagestats

import (
	"fmt"
	"regexp"

	"github.com/br0-space/bot/internal/db"
	"github.com/br0-space/bot/internal/matcher/abstract"
	"github.com/br0-space/bot/internal/matcher/registry"
	"github.com/br0-space/bot/internal/telegram"
)

// Each matcher extends the abstract matcher
type Matcher struct {
	abstract.Matcher
}

// Return the identifier of this matcher for use in logging
func (m Matcher) Identifier() string {
	return "messagestats"
}

// This is a command matcher and generates a help item
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{{
		Command:     "words",
		Description: "Zeigt alle dem Bot bekannten User an, sortiert nach der Anzahl ihrer bisher geschriebenen Worte",
	}}
}

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
	// Write stats on each post
	db.InsertMessageStats(requestMessage)

	// Check if text starts with /stats and if not, ignore it
	if doesMatch := m.doesMatch(requestMessage.Text); !doesMatch {
		return nil
	}

	records := db.GetWordCounts()

	return m.sendResponse(requestMessage, records)
}

// Check if a text starts with /stats
func (m Matcher) doesMatch(text string) bool {
	// Check if message starts with /choose
	match, _ := regexp.MatchString(`^/words(@|\s|$)`, text)

	return match
}

func (m Matcher) sendResponse(requestMessage telegram.RequestMessage, records []db.WordCount) error {
	responseText := "```"

	// Add one line per record
	for _, record := range records {
		responseText = responseText + fmt.Sprintf("\n%6d | %s", record.Words, record.Username)
	}

	responseText = responseText + "```"

	responseMessage := telegram.Message{
		Text:      responseText,
		ParseMode: "Markdown",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}
