package fortune

import (
	"fmt"
	"github.com/br0-space/bot/internal/telegram"
	"regexp"
	"strings"
)

const (
	typeText          Type   = "text"
	typeQuote         Type   = "quote"
	typeQuotePattern  string = `(?s)^(.+?)\n\n-- (.+)$`
	typeQuoteTemplate string = "%s\n\n_*\\-\\- %s*_"
	lineQuotePattern  string = `^(.+?): (.+)$`
	lineQuoteTemplate string = "*%s*: %s"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Type string

func getType(text string) Type {
	switch {
	case regexp.MustCompile(typeQuotePattern).MatchString(text):
		return typeQuote
	default:
		return typeText
	}
}

func (t Type) getFortune(text string) Fortune {
	text = strings.TrimSpace(text)

	var lines []string
	var source *string

	switch t {
	case typeText:
		lines = strings.Split(text, "\n")
	case typeQuote:
		matches := regexp.MustCompile(typeQuotePattern).FindStringSubmatch(text)
		lines = strings.Split(matches[1], "\n")
		source = &matches[2]
	}

	return Fortune{
		_type:   t,
		content: lines,
		source:  source,
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Fortune struct {
	_type   Type
	file    string
	content []string
	source  *string
}

func MakeFortune(file string, text string) Fortune {
	fortune := getType(text).getFortune(text)
	fortune.file = file

	return fortune
}

func (f Fortune) GetFile() string {
	return f.file
}

func (f Fortune) ToMarkdown() string {
	switch f._type {
	case typeText:
		return formatLines(f.content)
	case typeQuote:
		return fmt.Sprintf(
			typeQuoteTemplate,
			formatLines(f.content),
			telegram.EscapeMarkdown(*f.source),
		)
	default:
		return "unknown fortune type"
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func formatLines(lines []string) string {
	res := ""

	for _, line := range lines {
		res += formatLine(line) + "\n"
	}

	return strings.TrimSpace(res)
}

func formatLine(line string) string {
	if regexp.MustCompile(lineQuotePattern).MatchString(line) {
		matches := regexp.MustCompile(lineQuotePattern).FindStringSubmatch(line)
		return fmt.Sprintf(
			lineQuoteTemplate,
			telegram.EscapeMarkdown(matches[1]),
			telegram.EscapeMarkdown(matches[2]),
		)
	}

	return telegram.EscapeMarkdown(line)
}
