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
