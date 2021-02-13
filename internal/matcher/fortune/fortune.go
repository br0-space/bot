package fortune

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/neovg/kmptnzbot/internal/fortune"
	"github.com/neovg/kmptnzbot/internal/matcher/abstract"
	"github.com/neovg/kmptnzbot/internal/telegram"
)

// Each matcher extends the abstract matcher
type Matcher struct {
	abstract.Matcher
}

type helpItem struct {
	command string
	description string
}

// Return the identifier of this matcher for use in logging
func (m Matcher) Identifier() string {
	return "fortune"
}

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
	// Check if text starts with /fortune and if not, ignore it
	if doesMatch := m.doesMatch(requestMessage.Text); !doesMatch {
		return nil
	}

	// Get the option
	option := m.getOption(requestMessage.Text)

	// If no option was chosen, send a random fortune cookie
	if len(option) == 0 {
		return m.sendRandomFortune(requestMessage)
	}

	// help
	if option == "help" {
		return m.sendHelpResponse(requestMessage)
	}

	// list
	if option == "list" {
		return m.sendListResponse(requestMessage)
	}

	if fortune.Exists(option) {
		return m.sendFortune(requestMessage, option)
	}

	// Choose one option and send the result
	return fmt.Errorf("could not find cookie file `%s`", option)
}

// Check if a text starts with /jn or /yn
func (m Matcher) doesMatch(text string) bool {
	match, _ := regexp.MatchString(`^/fortune(\s|$)`, text)

	return match
}

// Check if a text starts with /jn and return the text behind
func (m Matcher) getOption(text string) string {
	match, _ := regexp.MatchString(`^/fortune .+`, text)
	if !match {
		return ""
	}

	return text[9:]
}

func (m Matcher) sendRandomFortune(requestMessage telegram.RequestMessage) error {
	f, err := fortune.GetRandomFortune()
	if err != nil {
		return err
	}

	return m.sendFortuneResponse(requestMessage, m.formatFortune(f))
}

func (m Matcher) sendFortune(requestMessage telegram.RequestMessage, file string) error {
	f, err := fortune.GetFortune(file)
	if err != nil {
		return err
	}

	return m.sendFortuneResponse(requestMessage, m.formatFortune(f))
}

func (m Matcher) sendFortuneResponse(requestMessage telegram.RequestMessage, text string) error {
	responseText := text

	responseMessage := telegram.Message{
		Text:             responseText,
		ReplyToMessageID: requestMessage.ID,
		ParseMode:        "HTML",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}

func (m Matcher) sendHelpResponse(requestMessage telegram.RequestMessage) error {
	helpItems := make([]helpItem, 0, 0)
	helpItems = append(helpItems, helpItem{
		command:     "",
		description: "Zeigt ein zufällig ausgewähltes Fortune-Cookie an",
	})
	helpItems = append(helpItems, helpItem{
		command:     "help",
		description: "Zeigt diese Hilfe an",
	})
	helpItems = append(helpItems, helpItem{
		command:     "list",
		description: "Zeigt eine Liste verfügbarer Cookie-Dateien an",
	})
	helpItems = append(helpItems, helpItem{
		command:     "{cookie-datei}",
		description: "Zeigt ein Fortune-Cookie aus der angegebenen Datei an",
	})

	responseText := ""

	for _, helpItem := range helpItems {
		responseText = responseText + fmt.Sprintf("\n`/fortune %s`\n   _%s_\n", helpItem.command, helpItem.description)
	}

	responseMessage := telegram.Message{
		Text:             responseText,
		ReplyToMessageID: requestMessage.ID,
		ParseMode:        "Markdown",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}

func (m Matcher) sendListResponse(requestMessage telegram.RequestMessage) error {
	files := fortune.GetList()

	responseText := "```"

	// Add one line per file
	for _, file := range files {
		responseText = responseText + fmt.Sprintf("\n%s", file)
	}

	responseText = responseText + "```"

	responseMessage := telegram.Message{
		Text:             responseText,
		ReplyToMessageID: requestMessage.ID,
		ParseMode:        "Markdown",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}

func (m Matcher) formatFortune(text string) string {
	lines := strings.Split(text, "\n")

	if m.looksLikeQuote(lines) {
		return m.formatQuoteFortune(lines)
	}

	return text
}

func (m Matcher) looksLikeQuote(lines []string) bool {
	isQuote, _ := regexp.MatchString(`^\s*-- `, lines[len(lines)-1])

	return isQuote
}

func (m Matcher) formatQuoteFortune(lines []string) string {
	text := ""
	for i, line := range lines {
		if i < len(lines) - 1 {
			if len(text) == 0 || text[len(text)-1:] == "\n" {
				text = text + "<i>"
			}
			text = text + line
			if len(line) >= 50 {
				if i == len(lines) - 2 {
					text = text + "</i>\n"
				} else {
					text = text + " "
				}
			} else {
				text = text + "</i>\n"
			}
		}
	}

	return text + "\n" + lines[len(lines)-1]
}
