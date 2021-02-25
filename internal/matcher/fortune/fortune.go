package fortune

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kmptnz/bot/internal/fortune"
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
	return "fortune"
}

// This is a command matcher and generates a help item
func (m Matcher) GetHelpItems() []registry.HelpItem {
	return []registry.HelpItem{{
		Command:     "fortune",
		Description: "Zeigt ein Fortune Cookie an. Siehe `/fortune help` für mehr Optionen.",
	}}
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

	// Show help
	if option == "help" {
		return m.sendHelpResponse(requestMessage)
	}

	// Show list of fortune files
	if option == "list" {
		return m.sendListResponse(requestMessage)
	}

	// Show a random fortune from a specific file
	if fortune.Exists(option) {
		return m.sendFortune(requestMessage, option)
	}

	// Choose one option and send the result
	return fmt.Errorf("could not find cookie file `%s`", option)
}

// Check if a text starts with /fortune
func (m Matcher) doesMatch(text string) bool {
	match, _ := regexp.MatchString(`^/fortune(@|\s|$)`, text)

	return match
}

// Check if a text starts with /fortune and return the text behind
func (m Matcher) getOption(text string) string {
	// Initialize the regular expression
	r := regexp.MustCompile(`\s\S+`)

	// Find the first word
	option := r.FindString(text)

	// Return a trimmed version
	return strings.TrimSpace(option)
}

// Send a random fortune from a random file
func (m Matcher) sendRandomFortune(requestMessage telegram.RequestMessage) error {
	f, err := fortune.GetRandomFortune()
	if err != nil {
		return err
	}

	return m.sendFortuneResponse(requestMessage, m.formatFortune(f))
}

// Send a random fortune from a specific file
func (m Matcher) sendFortune(requestMessage telegram.RequestMessage, file string) error {
	f, err := fortune.GetFortune(file)
	if err != nil {
		return err
	}

	return m.sendFortuneResponse(requestMessage, m.formatFortune(f))
}

// Send the formatted text of a fortune cookie
func (m Matcher) sendFortuneResponse(requestMessage telegram.RequestMessage, text string) error {
	responseText := text

	responseMessage := telegram.Message{
		Text:      responseText,
		ParseMode: "HTML",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}

// Send a list of available /fortune commands
func (m Matcher) sendHelpResponse(requestMessage telegram.RequestMessage) error {
	helpItems := make([]registry.HelpItem, 0)
	helpItems = append(helpItems, registry.HelpItem{
		Command:     "",
		Description: "Zeigt ein zufällig ausgewähltes Fortune-Cookie an",
	})
	helpItems = append(helpItems, registry.HelpItem{
		Command:     "help",
		Description: "Zeigt diese Hilfe an",
	})
	helpItems = append(helpItems, registry.HelpItem{
		Command:     "list",
		Description: "Zeigt eine Liste verfügbarer Cookie-Dateien an",
	})
	helpItems = append(helpItems, registry.HelpItem{
		Command:     "{cookie-datei}",
		Description: "Zeigt ein Fortune-Cookie aus der angegebenen Datei an",
	})

	responseText := ""

	for _, helpItem := range helpItems {
		responseText = responseText + fmt.Sprintf("\n`/fortune %s`\n   _%s_\n", helpItem.Command, helpItem.Description)
	}

	responseMessage := telegram.Message{
		Text:      responseText,
		ParseMode: "Markdown",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}

// Send a list of available fortune files
func (m Matcher) sendListResponse(requestMessage telegram.RequestMessage) error {
	files := fortune.GetList()

	responseText := "```"

	// Add one line per file
	for _, file := range files {
		responseText = responseText + fmt.Sprintf("\n%s", file)
	}

	responseText = responseText + "```"

	responseMessage := telegram.Message{
		Text:      responseText,
		ParseMode: "Markdown",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}

// Try to format a fortune cookie in a way which is nice to read in Telegram
func (m Matcher) formatFortune(text string) string {
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")

	lines := strings.Split(text, "\n")

	if m.looksLikeQuote(lines) {
		return m.formatQuoteFortune(lines)
	}

	return text
}

// Check if a bunch of lines look like a quote (last line begins with `--`)
func (m Matcher) looksLikeQuote(lines []string) bool {
	isQuote, _ := regexp.MatchString(`^\s*-- `, lines[len(lines)-1])

	return isQuote
}

// Format a bunch of lines as quote
// Tries to remove hard wraps while preserving intentional line breaks
// The quote itself will become italic, the last line with the author info will remain normal text
func (m Matcher) formatQuoteFortune(lines []string) string {
	text := ""
	for i, line := range lines {
		// Ignore the last line (with the author of the quote)
		if i < len(lines)-1 {
			// Every line in the formatted quote starts with <i>
			if len(text) == 0 || text[len(text)-1:] == "\n" {
				text = text + "<i>"
			}
			// Append the next line to the new text
			text = text + line
			// We view line breaks at the end of lines of 50 characters or more as  hard wraps (like here:)),
			// while we view lines with less than 50 characters as end of a paragraph
			if len(line) >= 50 {
				// Long lines will be directly followed by the next line, the line break is removed
				if i == len(lines)-2 {
					// The last line of a quote always ends with </i> and a line break, regardless of its length
					text = text + "</i>\n"
				} else {
					// Normally long line just get the next line appended after a white space
					text = text + " "
				}
			} else {
				// Short lines probably end a paragraph, hence they are followed by </i> and a line break
				text = text + "</i>\n"
			}
		}
	}

	// Return the quote followed by the unformatted author info
	return text + "\n" + lines[len(lines)-1]
}
