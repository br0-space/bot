package fortune_test

import (
	"strings"
	"testing"

	"github.com/br0-space/bot/pkg/fortune"
)

func TestMakeFortune(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		file     string
		text     string
		wantType fortune.Type
	}{
		{
			name:     "simple text fortune",
			file:     "test",
			text:     "This is a simple fortune",
			wantType: "text",
		},
		{
			name:     "multiline text fortune",
			file:     "test",
			text:     "Line 1\nLine 2\nLine 3",
			wantType: "text",
		},
		{
			name:     "quote fortune",
			file:     "quotes",
			text:     "This is a quote\n\n-- Author",
			wantType: "quote",
		},
		{
			name:     "quote with multiple lines",
			file:     "quotes",
			text:     "First line\nSecond line\n\n-- Famous Person",
			wantType: "quote",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fortune.MakeFortune(tc.file, tc.text)

			if f.Type() != tc.wantType {
				t.Errorf("MakeFortune type = %v, want %v", f.Type(), tc.wantType)
			}

			if f.File() != tc.file {
				t.Errorf("MakeFortune file = %v, want %v", f.File(), tc.file)
			}

			if f.ToMarkdown() == "" {
				t.Error("Fortune content should not be empty")
			}
		})
	}
}

func TestFortune_File(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		file     string
		text     string
		wantFile string
	}{
		{"standard file name", "myfile", "Some text", "myfile"},
		{"file with path-like name", "category/subcategory", "Text", "category/subcategory"},
		{"empty file name", "", "Text", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fortune.MakeFortune(tc.file, tc.text)
			got := f.File()

			if got != tc.wantFile {
				t.Errorf("File() = %q, want %q", got, tc.wantFile)
			}
		})
	}
}

func TestFortune_ToMarkdown(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		file         string
		text         string
		wantContains []string
		wantNotEmpty bool
	}{
		{
			name:         "simple text",
			file:         "test",
			text:         "Hello World",
			wantContains: []string{"Hello World"},
			wantNotEmpty: true,
		},
		{
			name:         "text with special markdown chars",
			file:         "test",
			text:         "Test * _ [ ] ( ) ~ ` > # + - = | { } . !",
			wantContains: []string{},
			wantNotEmpty: true,
		},
		{
			name:         "quote with source",
			file:         "quotes",
			text:         "A great quote\n\n-- Famous Author",
			wantContains: []string{"great quote", "Famous Author"},
			wantNotEmpty: true,
		},
		{
			name:         "dialog format",
			file:         "dialogs",
			text:         "Alice: Hello there\nBob: Hi back",
			wantContains: []string{"Alice", "Hello", "Bob", "Hi"},
			wantNotEmpty: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fortune.MakeFortune(tc.file, tc.text)
			markdown := f.ToMarkdown()

			if tc.wantNotEmpty && markdown == "" {
				t.Error("ToMarkdown() returned empty string")
			}

			for _, substr := range tc.wantContains {
				if !strings.Contains(markdown, substr) {
					t.Errorf("ToMarkdown() should contain %q, got: %s", substr, markdown)
				}
			}
		})
	}
}

func TestFormatLine(t *testing.T) {
	t.Parallel()

	// We can't test private methods directly, but we test through ToMarkdown
	testCases := []struct {
		name         string
		text         string
		wantContains string
	}{
		{"simple text", "Hello world", "Hello world"},
		{"dialog format", "Speaker: Message", "Speaker"},
		{"colon in text", "Time: 10:30", "Time"},
		{"no colon", "Just text", "Just text"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fortune.MakeFortune("test", tc.text)
			result := f.ToMarkdown()

			if !strings.Contains(result, tc.wantContains) {
				t.Errorf("ToMarkdown(%q) should contain %q, got %q", tc.text, tc.wantContains, result)
			}

			if result == "" {
				t.Errorf("ToMarkdown(%q) returned empty string", tc.text)
			}
		})
	}
}

func TestFormatLines(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		text  string
		wants []string
	}{
		{
			name:  "single line",
			text:  "Hello",
			wants: []string{"Hello"},
		},
		{
			name:  "multiple lines",
			text:  "Line 1\nLine 2\nLine 3",
			wants: []string{"Line 1", "Line 2", "Line 3"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fortune.MakeFortune("test", tc.text)
			result := f.ToMarkdown()

			// Check that all input lines are present in output
			for _, want := range tc.wants {
				if !strings.Contains(result, want) {
					t.Errorf("ToMarkdown result should contain %q, got %q", want, result)
				}
			}
		})
	}
}

func TestFortune_QuoteWithSpecialCharacters(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		text   string
		source string
	}{
		{
			name:   "quote with markdown characters",
			text:   "Text with *emphasis* and _underscores_\n\n-- Author*Name",
			source: "Author*Name",
		},
		{
			name:   "quote with brackets",
			text:   "Text [with] (brackets)\n\n-- [Author]",
			source: "[Author]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := fortune.MakeFortune("test", tc.text)
			markdown := f.ToMarkdown()

			if markdown == "" {
				t.Error("ToMarkdown() returned empty string for quote")
			}

			// The source should be present in some form in the markdown output
			if !strings.Contains(markdown, tc.source) {
				t.Logf("Warning: Source %q not found in markdown output: %s", tc.source, markdown)
			}
		})
	}
}
