package janein_test

import (
	"strings"
	"testing"

	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/pkg/matchers/janein"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expected = struct {
	insult string
	yes    string
	no     string
}{
	insult: `Ob du behindert bist hab ich gefragt?\! 🤪`,
	yes:    `👍 *Ja*, du solltest *foo\* bar\_*\!`,
	no:     `👎 *Nein*, du solltest nicht *foo\* bar\_*\!`,
}

var expectedReply = []telegramclient.MessageStruct{{
	ChatID:                0,
	ReplyToMessageID:      123,
	Text:                  "",
	Photo:                 "",
	Caption:               "",
	ParseMode:             "MarkdownV2",
	DisableWebPagePreview: false,
	DisableNotification:   false,
}}

var tests = []struct {
	in              string
	expectedReplies []telegramclient.MessageStruct
}{
	{"", nil},
	{"foobar", nil},
	{"jn", nil},
	{"/jnx", nil},
	{" /jn", nil},
	{"/jn", expectedReply},
	{"/jn@bot", expectedReply},
	{"/yn", expectedReply},
	{"/yn@bot", expectedReply},
	{"/jn foo* bar_", expectedReply},
	{"/jn@bot foo* bar_", expectedReply},
	{"/yn foo* bar_", expectedReply},
	{"/yn@bot foo* bar_", expectedReply},
}

func provideMatcher() janein.Matcher {
	return janein.MakeMatcher()
}

func newTestMessage(text string) telegramclient.WebhookMessageStruct {
	return telegramclient.TestWebhookMessage(text)
}

func TestMatcher_DoesMatch(t *testing.T) {
	t.Parallel()

	for _, tt := range tests {
		doesMatch := provideMatcher().DoesMatch(newTestMessage(tt.in))
		assert.Equal(t, tt.expectedReplies != nil, doesMatch, tt.in)
	}
}

func TestMatcher_Process(t *testing.T) {
	t.Parallel()

	for _, tt := range tests {
		replies, err := provideMatcher().Process(newTestMessage(tt.in))
		if tt.expectedReplies == nil {
			require.Error(t, err, tt.in)
			assert.Nil(t, replies, tt.in)
		} else {
			require.NoError(t, err, tt.in)
			assert.NotNil(t, replies, tt.in)
			assert.Len(t, replies, 1, tt.in)

			expectedReplies := tt.expectedReplies

			switch {
			case strings.Contains((replies)[0].Text, "behindert"):
				expectedReplies[0].Text = expected.insult
			case strings.Contains((replies)[0].Text, "Ja"):
				expectedReplies[0].Text = expected.yes
			case strings.Contains((replies)[0].Text, "Nein"):
				expectedReplies[0].Text = expected.no
			}

			assert.Equal(t, expectedReplies, replies, tt.in)
		}
	}
}
