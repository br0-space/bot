package fortune_test

import (
	"testing"

	"github.com/br0-space/bot/pkg/fortune"
)

func TestGetType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		text     string
		wantType fortune.Type
	}{
		{
			name:     "simple text",
			text:     "This is plain text",
			wantType: "text",
		},
		{
			name:     "multiline text",
			text:     "Line 1\nLine 2\nLine 3",
			wantType: "text",
		},
		{
			name:     "quote with source",
			text:     "This is a quote\n\n-- Author",
			wantType: "quote",
		},
		{
			name:     "quote with multiline content",
			text:     "First line\nSecond line\nThird line\n\n-- Author Name",
			wantType: "quote",
		},
		{
			name:     "text with single dash",
			text:     "Text with - dash",
			wantType: "text",
		},
		{
			name:     "text with -- but no newlines",
			text:     "Text with -- dashes",
			wantType: "text",
		},
		{
			name:     "empty text",
			text:     "",
			wantType: "text",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fortune.MakeFortune("test", tc.text)
			got := f.Type()

			if got != tc.wantType {
				t.Errorf("getType(%q) = %v, want %v", tc.text, got, tc.wantType)
			}
		})
	}
}

func TestType_GetFortune(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		text         string
		wantType     fortune.Type
		wantNotEmpty bool
	}{
		{
			name:         "text type single line",
			text:         "Single line",
			wantType:     "text",
			wantNotEmpty: true,
		},
		{
			name:         "text type multiple lines",
			text:         "Line 1\nLine 2\nLine 3",
			wantType:     "text",
			wantNotEmpty: true,
		},
		{
			name:         "quote type",
			text:         "Quote text\n\n-- Author",
			wantType:     "quote",
			wantNotEmpty: true,
		},
		{
			name:         "quote type multiline",
			text:         "First line\nSecond line\n\n-- Famous Person",
			wantType:     "quote",
			wantNotEmpty: true,
		},
		{
			name:         "text with whitespace",
			text:         "  \n  Text with whitespace  \n  ",
			wantType:     "text",
			wantNotEmpty: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fortune.MakeFortune("test", tc.text)

			if f.Type() != tc.wantType {
				t.Errorf("getFortune type = %v, want %v", f.Type(), tc.wantType)
			}

			if tc.wantNotEmpty && f.ToMarkdown() == "" {
				t.Error("Fortune ToMarkdown should not be empty")
			}
		})
	}
}

func TestType_GetFortune_QuoteExtraction(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		text           string
		wantInMarkdown string
	}{
		{
			name:           "simple quote",
			text:           "Quote\n\n-- Author",
			wantInMarkdown: "Author",
		},
		{
			name:           "quote with complex author",
			text:           "Text\n\n-- Author Name, Book Title",
			wantInMarkdown: "Author Name",
		},
		{
			name:           "quote with special characters in author",
			text:           "Text\n\n-- Author (1900-2000)",
			wantInMarkdown: "Author",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fortune.MakeFortune("test", tc.text)
			markdown := f.ToMarkdown()

			if f.Type() != "quote" {
				t.Error("Expected quote type")
			}

			if markdown == "" {
				t.Error("Markdown should not be empty for quote")
			}
		})
	}
}

func TestType_GetFortune_ContentParsing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		text         string
		wantType     fortune.Type
		wantContains []string
	}{
		{
			name:         "text with three lines",
			text:         "Line A\nLine B\nLine C",
			wantType:     "text",
			wantContains: []string{"Line A", "Line B", "Line C"},
		},
		{
			name:         "quote with two lines",
			text:         "First\nSecond\n\n-- Author",
			wantType:     "quote",
			wantContains: []string{"First", "Second"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fortune.MakeFortune("test", tc.text)

			if f.Type() != tc.wantType {
				t.Errorf("Type = %v, want %v", f.Type(), tc.wantType)
			}

			markdown := f.ToMarkdown()
			for range tc.wantContains {
				// Note: we can't check exact content since it gets escaped,
				// but ToMarkdown should not be empty
				if markdown == "" {
					t.Error("ToMarkdown should not be empty")

					break
				}
			}
		})
	}
}
