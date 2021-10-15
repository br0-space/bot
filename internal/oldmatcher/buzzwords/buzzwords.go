package buzzwords

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/br0-space/bot/internal/oldconfig"
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
	return "buzzwords"
}

// This matcher is no command and generates no help items
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{}
}

// Process a message received from telegram
func (m Matcher) ProcessRequestMessage(requestMessage oldtelegram.RequestMessage) error {
	// Create a wait group for synchronization
	var waitGroup sync.WaitGroup

	// We need to wait until all buzzwords are processed
	waitGroup.Add(len(oldconfig.Cfg.BuzzwordsMatcher))

	// Launch a goroutine for each token
	for _, buzzword := range oldconfig.Cfg.BuzzwordsMatcher {
		go func(buzzword oldconfig.BuzzwordsMatcher) {
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
func (m Matcher) processBuzzword(requestMessage oldtelegram.RequestMessage, buzzword oldconfig.BuzzwordsMatcher) error {
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

	// Send message to telegram
	return m.sendResponse(requestMessage, buzzword, newValue)
}

// Send a message with the new value to telegram
func (m Matcher) sendResponse(requestMessage oldtelegram.RequestMessage, buzzword oldconfig.BuzzwordsMatcher, newValue int) error {
	responseMessage := oldtelegram.Message{
		Text:      fmt.Sprintf(buzzword.Reply, newValue),
		ParseMode: "Markdown",
	}

	return oldtelegram.SendMessage(requestMessage, responseMessage)
}
