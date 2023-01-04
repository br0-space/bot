package plusplus_test

import (
	"github.com/br0-space/bot/container"
	"github.com/br0-space/bot/internal/matcher/plusplus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var parseTests = []struct {
	in       []string
	expected []plusplus.Token
}{
	{[]string{}, []plusplus.Token{}},
	{[]string{"foo++"}, []plusplus.Token{{"foo", 1}}},
	{[]string{"foo--"}, []plusplus.Token{{"foo", -1}}},
	{[]string{"foo+-"}, []plusplus.Token{{"foo", 0}}},
	{[]string{"foo—"}, []plusplus.Token{{"foo", -1}}},
	{[]string{"foo++", "foo++"}, []plusplus.Token{{"foo", 2}}},
	{[]string{"foo++", "foo--"}, []plusplus.Token{{"foo", 0}}},
	{[]string{"foo++", "bar--"}, []plusplus.Token{{"foo", 1}, {"bar", -1}}},
	{[]string{"foo+++++"}, []plusplus.Token{{"foo", 4}}},
	{[]string{"foo-----"}, []plusplus.Token{{"foo", -4}}},
}

func provideMatcher() plusplus.Matcher {
	return plusplus.NewMatcher(
		container.ProvideLogger(),
		nil,
	)
}

func TestMatcher_ParseTokens(t *testing.T) {
	t.Parallel()

	for _, tt := range parseTests {
		tokens := provideMatcher().ParseTokens(tt.in)
		assert.Equal(t, tt.expected, tokens, tt.in)
	}
}
