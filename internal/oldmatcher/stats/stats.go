package stats

import (
	"fmt"
	"regexp"

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
	return "stats"
}

// This is a command matcher and generates a help item
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{{
		Command:     "stats",
		Description: "Zeigt alle dem Bot bekannten User an, sortiert nach der Anzahl ihrer bisherigen Posts",
	}, {
		Command:     "words",
		Description: "Zeigt alle dem Bot bekannten User an, sortiert nach der Anzahl ihrer bisher geschriebenen Worte",
	}}
}

// Process a message received from telegram
func (m Matcher) ProcessRequestMessage(requestMessage oldtelegram.RequestMessage) error {
	// Write stats on each post
	db.UpdateStats(requestMessage.From)

	// Check if text starts with /stats and if not, ignore it
	if doesMatch := m.doesMatch(requestMessage.Text); !doesMatch {
		return nil
	}

	records := db.FindStatsTop()

	return m.sendResponse(requestMessage, records)
}

// Check if a text starts with /stats
func (m Matcher) doesMatch(text string) bool {
	// Check if message starts with /choose
	match, _ := regexp.MatchString(`^/stats(@|\s|$)`, text)

	return match
}

func (m Matcher) sendResponse(requestMessage oldtelegram.RequestMessage, records []db.Stats) error {
	responseText := "```"

	// Add one line per record
	for _, record := range records {
		responseText = responseText + fmt.Sprintf("\n%6d | %s", record.Posts, record.Username)
	}

	responseText = responseText + "```"

	responseMessage := oldtelegram.Message{
		Text:      responseText,
		ParseMode: "Markdown",
	}

	return oldtelegram.SendMessage(requestMessage, responseMessage)
}
