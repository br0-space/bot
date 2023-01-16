package fortune

import (
	"fmt"
	"github.com/br0-space/bot/pkg/telegram"
	"regexp"
	"strings"
)

const (
	typeQuoteTemplate string = "%s\n\n_*\\-\\- %s*_"
	lineQuotePattern  string = `^(.+?): (.+)$`
	lineQuoteTemplate string = "*%s*: %s"
)

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

func (f Fortune) File() string {
	return f.file
}

func (f Fortune) ToMarkdown() string {
	switch f._type {
	case typeText:
		return f.formatLines(f.content)
	case typeQuote:
		return fmt.Sprintf(
			typeQuoteTemplate,
			f.formatLines(f.content),
			telegram.EscapeMarkdown(*f.source),
		)
	default:
		return "unknown fortune type"
	}
}

func (f Fortune) formatLines(lines []string) string {
	res := ""

	for _, line := range lines {
		res += f.formatLine(line) + "\n"
	}

	return strings.TrimSpace(res)
}

func (f Fortune) formatLine(line string) string {
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
