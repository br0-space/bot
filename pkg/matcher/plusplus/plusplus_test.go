package plusplus_test

import (
	"fmt"
	"github.com/br0-space/bot/pkg/matcher/plusplus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var parseTests = []struct {
	in       []string
	expected []plusplus.Token
	err      error
}{
	{[]string{}, []plusplus.Token{}, nil},
	{[]string{"foo++"}, []plusplus.Token{{"foo", 1}}, nil},
	{[]string{"foo+++"}, []plusplus.Token{{"foo", 2}}, nil},
	{[]string{"foo--"}, []plusplus.Token{{"foo", -1}}, nil},
	{[]string{"foo---"}, []plusplus.Token{{"foo", -2}}, nil},
	{[]string{"foo+-"}, []plusplus.Token{{"foo", 0}}, nil},
	{[]string{"foo+--"}, []plusplus.Token{{"foo+", -1}}, nil},
	{[]string{"foo-+"}, []plusplus.Token{{"foo", 0}}, nil},
	{[]string{"foo-++"}, []plusplus.Token{{"foo-", 1}}, nil},
	{[]string{"foo—"}, []plusplus.Token{{"foo", -1}}, nil},
	{[]string{"foo++", "foo++"}, []plusplus.Token{{"foo", 2}}, nil},
	{[]string{"foo++", "foo--"}, []plusplus.Token{{"foo", 0}}, nil},
	{[]string{"foo++", "bar--"}, []plusplus.Token{{"foo", 1}, {"bar", -1}}, nil},
	{[]string{"foo+bar++"}, []plusplus.Token{{"foo+bar", 1}}, nil},
	{[]string{"foo++bar++"}, []plusplus.Token{{"foo++bar", 1}}, nil},
	{[]string{"foo"}, nil, fmt.Errorf(`unable to find mode in match "foo"`)},
	{[]string{"++"}, nil, fmt.Errorf(`unable to find name in match "++"`)},
}

func TestParseTokens(t *testing.T) {
	t.Parallel()

	for _, tt := range parseTests {
		tokens, err := plusplus.GetTokens(tt.in)
		assert.Equal(t, tt.expected, tokens, tt.in)
		assert.Equal(t, tt.err, err, tt.in)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var getTokenIncrementTests = []struct {
	in       string
	expected int
	err      error
}{
	{"++", 1, nil},
	{"+++", 2, nil},
	{"++++", 3, nil},
	{"--", -1, nil},
	{"---", -2, nil},
	{"----", -3, nil},
	{"+-", 0, nil},
	{"-+", 0, nil},
	{"—", -1, nil},
	{"foo", 0, fmt.Errorf(`unable to get increment value from mode "foo"`)},
	{"+--", 0, fmt.Errorf(`unable to get increment value from mode "+--"`)},
	{"-++", 0, fmt.Errorf(`unable to get increment value from mode "-++"`)},
}

func TestGetTokenIncrement(t *testing.T) {
	t.Parallel()

	for _, tt := range getTokenIncrementTests {
		increment, err := plusplus.GetTokenIncrement(tt.in)
		assert.Equal(t, tt.expected, increment, tt.in)
		assert.Equal(t, tt.err, err, tt.in)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//func provideMatcher() plusplus.Matcher {
//	return plusplus.NewMatcher(
//		container.ProvideLogger(),
//		nil,
//	)
//}
