package plusplus_test

import (
	"fmt"
	"github.com/br0-space/bot/internal/matcher/plusplus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var parseTests = []struct {
	in       []string
	err      error
	expected []plusplus.Token
}{
	{[]string{}, nil, []plusplus.Token{}},
	{[]string{"foo++"}, nil, []plusplus.Token{{"foo", 1}}},
	{[]string{"foo+++"}, nil, []plusplus.Token{{"foo", 2}}},
	{[]string{"foo--"}, nil, []plusplus.Token{{"foo", -1}}},
	{[]string{"foo---"}, nil, []plusplus.Token{{"foo", -2}}},
	{[]string{"foo+-"}, nil, []plusplus.Token{{"foo", 0}}},
	{[]string{"foo+--"}, nil, []plusplus.Token{{"foo+", -1}}},
	{[]string{"foo-+"}, nil, []plusplus.Token{{"foo", 0}}},
	{[]string{"foo-++"}, nil, []plusplus.Token{{"foo-", 1}}},
	{[]string{"foo—"}, nil, []plusplus.Token{{"foo", -1}}},
	{[]string{"foo++", "foo++"}, nil, []plusplus.Token{{"foo", 2}}},
	{[]string{"foo++", "foo--"}, nil, []plusplus.Token{{"foo", 0}}},
	{[]string{"foo++", "bar--"}, nil, []plusplus.Token{{"foo", 1}, {"bar", -1}}},
	{[]string{"foo+bar++"}, nil, []plusplus.Token{{"foo+bar", 1}}},
	{[]string{"foo++bar++"}, nil, []plusplus.Token{{"foo++bar", 1}}},
	{[]string{"foo"}, fmt.Errorf(`unable to find mode in match "foo"`), nil},
	{[]string{"++"}, fmt.Errorf(`unable to find name in match "++"`), nil},
}

func TestParseTokens(t *testing.T) {
	t.Parallel()

	for _, tt := range parseTests {
		err, tokens := plusplus.GetTokens(tt.in)
		assert.Equal(t, tt.err, err, tt.in)
		assert.Equal(t, tt.expected, tokens, tt.in)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var getTokenIncrementTests = []struct {
	in       string
	err      error
	expected int
}{
	{"++", nil, 1},
	{"+++", nil, 2},
	{"++++", nil, 3},
	{"--", nil, -1},
	{"---", nil, -2},
	{"----", nil, -3},
	{"+-", nil, 0},
	{"-+", nil, 0},
	{"—", nil, -1},
	{"foo", fmt.Errorf(`unable to get increment value from mode "foo"`), 0},
	{"+--", fmt.Errorf(`unable to get increment value from mode "+--"`), 0},
	{"-++", fmt.Errorf(`unable to get increment value from mode "-++"`), 0},
}

func TestGetTokenIncrement(t *testing.T) {
	t.Parallel()

	for _, tt := range getTokenIncrementTests {
		err, increment := plusplus.GetTokenIncrement(tt.in)
		assert.Equal(t, tt.err, err, tt.in)
		assert.Equal(t, tt.expected, increment, tt.in)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//func provideMatcher() plusplus.Matcher {
//	return plusplus.NewMatcher(
//		container.ProvideLogger(),
//		nil,
//	)
//}
