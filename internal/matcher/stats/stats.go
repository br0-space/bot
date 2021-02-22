package stats

import (
	"fmt"
	"regexp"

	"github.com/kmptnz/bot/internal/db"
	"github.com/kmptnz/bot/internal/matcher/abstract"
	"github.com/kmptnz/bot/internal/matcher/registry"
	"github.com/kmptnz/bot/internal/telegram"
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

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
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

func (m Matcher) sendResponse(requestMessage telegram.RequestMessage, records []db.Stats) error {
	responseText := "```"

	// Add one line per record
	for _, record := range records {
		responseText = responseText + fmt.Sprintf("\n%6d | %s", record.Posts, record.Username)
	}

	responseText = responseText + "```"

	responseMessage := telegram.Message{
		Text:      responseText,
		ParseMode: "Markdown",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}
