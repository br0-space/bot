package fortune

import (
	"regexp"
	"strings"
)

const (
	typeText         Type   = "text"
	typeQuote        Type   = "quote"
	typeQuotePattern string = `(?s)^(.+?)\n\n-- (.+)$`
)

// Type represents the format type of a fortune entry (text or quote).
type Type string

// getType determines the type of fortune based on the text content.
// It returns typeQuote if the text matches the quote pattern, otherwise typeText.
func getType(text string) Type {
	switch {
	case regexp.MustCompile(typeQuotePattern).MatchString(text):
		return typeQuote
	default:
		return typeText
	}
}

// getFortune parses raw text into a Fortune struct based on the type.
// For quotes, it extracts the content and source attribution.
// For text, it simply splits the content into lines.
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
		file:    "",
		content: lines,
		source:  source,
	}
}
