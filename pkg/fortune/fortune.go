package fortune

import (
	"fmt"
	"regexp"
	"strings"

	telegramclient "github.com/br0-space/bot-telegramclient"
)

const (
	typeQuoteTemplate string = "%s\n\n_*\\-\\- %s*_"
	lineQuotePattern  string = `^(.+?): (.+)$`
	lineQuoteTemplate string = "*%s*: %s"
)

// Fortune represents a single fortune entry with its type, content, and source.
// It can be formatted as markdown text for display.
type Fortune struct {
	_type   Type
	file    string
	content []string
	source  *string
}

// MakeFortune creates a Fortune instance from raw text. It automatically detects
// the type of fortune (text or quote) and parses it accordingly.
func MakeFortune(file string, text string) Fortune {
	fortune := getType(text).getFortune(text)
	fortune.file = file

	return fortune
}

// File returns the name of the fortune file this fortune came from.
func (f Fortune) File() string {
	return f.file
}

// Type returns the type of the fortune (text or quote).
func (f Fortune) Type() Type {
	return f._type
}

// ToMarkdown converts the fortune to a markdown-formatted string suitable for display.
// Quotes are formatted with the source attribution, and special characters are escaped
// for Telegram markdown compatibility.
func (f Fortune) ToMarkdown() string {
	switch f._type {
	case typeText:
		return f.formatLines(f.content)
	case typeQuote:
		return fmt.Sprintf(
			typeQuoteTemplate,
			f.formatLines(f.content),
			telegramclient.EscapeMarkdown(*f.source),
		)
	default:
		return "unknown fortune type"
	}
}

// formatLines formats multiple lines of fortune content, applying line-specific formatting.
func (f Fortune) formatLines(lines []string) string {
	res := ""

	for _, line := range lines {
		res += f.formatLine(line) + "\n"
	}

	return strings.TrimSpace(res)
}

// formatLine formats a single line of text. If the line matches the pattern "Speaker: Text",
// it formats it as a dialog line with emphasis on the speaker.
func (f Fortune) formatLine(line string) string {
	if regexp.MustCompile(lineQuotePattern).MatchString(line) {
		matches := regexp.MustCompile(lineQuotePattern).FindStringSubmatch(line)

		return fmt.Sprintf(
			lineQuoteTemplate,
			telegramclient.EscapeMarkdown(matches[1]),
			telegramclient.EscapeMarkdown(matches[2]),
		)
	}

	return telegramclient.EscapeMarkdown(line)
}
