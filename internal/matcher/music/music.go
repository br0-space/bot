package music

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/br0-space/bot/internal/logger"
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
	return "music"
}

// This matcher is no command and generates no help items
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{}
}

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
	// Get matches
	matches := m.getMatches(requestMessage.TextOrCaption())

	// Process matches
	m.processMatches(requestMessage, matches)

	return nil
}

// Return a list of Spotify or Apple Music URLs contained in a text
func (m Matcher) getMatches(text string) []string {
	// Check if message starts with / and if yes, ignore it
	match, _ := regexp.MatchString(`^/`, text)
	if match {
		return make([]string, 0)
	}

	const spotifyUrlPattern = "https?:\\/\\/open.spotify.com\\/(album|track)\\/.+?"
	const appleMusicUrlPattern = "https?:\\/\\/music.apple.com\\/[a-z]{2}\\/album\\/.+?"
	urlPattern := fmt.Sprintf("(%s|%s)(\\s|$)", spotifyUrlPattern, appleMusicUrlPattern)

	// Initialize the regular expression
	r := regexp.MustCompile(urlPattern)

	// Find all occurrences
	matches := r.FindAllString(text, -1)

	// Trim matches to get rid of trailing spaces
	for i := range matches {
		matches[i] = strings.TrimSpace(matches[i])
	}

	return matches
}

func (m Matcher) processMatches(requestMessage telegram.RequestMessage, matches []string) {
	for _, match := range matches {
		err := m.processMatch(requestMessage, match)
		if err != nil {
			m.HandleError(requestMessage, m.Identifier(), err)
		}
	}
}

func (m Matcher) processMatch(requestMessage telegram.RequestMessage, match string) error {
	logger.Log.Info(match)

	songlinkEntry, err := GetSonglinkEntry(match)
	if err != nil {
		m.HandleError(requestMessage, m.Identifier(), err)
	}

	err = m.sendResponse(requestMessage, *songlinkEntry)

	return nil
}

func (m Matcher) sendResponse(requestMessage telegram.RequestMessage, songlinkEntry SonglinkEntry) error {
	responseText := fmt.Sprintf(
		"%s:\n\n%s\n(%s)\n\n",
		songlinkEntry.Type,
		songlinkEntry.Title,
		songlinkEntry.Artist,
	)

	// The entry may not be available at each platform, so just add existing links
	if songlinkEntry.SpotifyURL != "" {
		responseText += fmt.Sprintf("Spotify: %s\n\n", songlinkEntry.SpotifyURL)
	}
	if songlinkEntry.AppleMusicURL != "" {
		responseText += fmt.Sprintf("Apple Music: %s\n\n", songlinkEntry.AppleMusicURL)
	}
	if songlinkEntry.YoutubeURL != "" {
		responseText += fmt.Sprintf("YouTube: %s\n\n", songlinkEntry.YoutubeURL)
	}

	// Add link so Songlink page for convenience
	responseText += fmt.Sprintf("Weitere Links bei Songlink: %s", songlinkEntry.SonglinkURL)

	responseMessage := telegram.Message{
		Text: responseText,
		// ParseMode: "Markdown",
		DisableWebPagePreview: true,
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}
