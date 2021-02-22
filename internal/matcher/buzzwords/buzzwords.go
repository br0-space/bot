package buzzwords

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/kmptnz/bot/internal/config"
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
	return "buzzwords"
}

// This matcher is no command and generates no help items
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{}
}

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
	// Create a wait group for synchronization
	var waitGroup sync.WaitGroup

	// We need to wait until all buzzwords are processed
	waitGroup.Add(len(config.Cfg.BuzzwordsMatcher))

	// Launch a goroutine for each token
	for _, buzzword := range config.Cfg.BuzzwordsMatcher {
		go func(buzzword config.BuzzwordsMatcher) {
			defer waitGroup.Done()

			// Process the token
			err := m.processBuzzword(requestMessage, buzzword)
			if err != nil {
				m.HandleError(requestMessage, m.Identifier(), err)
			}
		}(buzzword)
	}

	// Wait until all tokens are processed
	waitGroup.Wait()

	return nil
}

// Check if the message contains a buzzword and process it
func (m Matcher) processBuzzword(requestMessage telegram.RequestMessage, buzzword config.BuzzwordsMatcher) error {
	// Check if message contains the trigger
	pattern := fmt.Sprintf("(?i)(^|\\s)%s(\\W|$)", buzzword.Trigger)
	match, _ := regexp.MatchString(pattern, requestMessage.Text)
	if !match {
		return nil
	}

	// Check if trigger is in plusplus
	pattern = fmt.Sprintf("(?i)(^|\\s)%s(\\+\\+|\\+-|--|—)", buzzword.Trigger)
	match, _ = regexp.MatchString(pattern, requestMessage.Text)
	if match {
		return nil
	}

	// Update the database
	newValue := db.IncrementPlusplus(buzzword.Trigger, 1)

	// Send message to Telegram
	return m.sendResponse(requestMessage, buzzword, newValue)
}

// Send a message with the new value to Telegram
func (m Matcher) sendResponse(requestMessage telegram.RequestMessage, buzzword config.BuzzwordsMatcher, newValue int) error {
	responseMessage := telegram.Message{
		Text:      fmt.Sprintf(buzzword.Reply, newValue),
		ParseMode: "Markdown",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}
