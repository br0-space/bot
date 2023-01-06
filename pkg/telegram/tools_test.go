package telegram_test

import (
	"github.com/br0-space/bot/pkg/telegram"
	"testing"
)

type escapeMarkdownTest struct {
	in  string
	out string
}

var escapeMarkdownTests = []escapeMarkdownTest{
	{"_", "\\_"},
	{"*", "\\*"},
	{"[", "\\["},
	{"]", "\\]"},
	{"(", "\\("},
	{")", "\\)"},
	{"~", "\\~"},
	{"`", "\\`"},
	{">", "\\>"},
	{"#", "\\#"},
	{"+", "\\+"},
	{"-", "\\-"},
	{"=", "\\="},
	{"|", "\\|"},
	{"{", "\\{"},
	{"}", "\\}"},
	{".", "\\."},
	{"!", "\\!"},
	{"\\", "\\\\"},
}

func TestEscapeMarkdown(t *testing.T) {
	t.Parallel()

	for _, test := range escapeMarkdownTests {
		if out := telegram.EscapeMarkdown(test.in); out != test.out {
			t.Errorf("EscapeMarkdown(%q) = %q, want %q", test.in, out, test.out)
		}
	}
}
