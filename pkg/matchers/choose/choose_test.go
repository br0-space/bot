package choose_test

import (
	"fmt"
	"strings"
	"testing"

	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/pkg/matchers/choose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expected = struct {
	insult  string
	success string
}{
	insult:  `Ob du behindert bist hab ich gefragt?\! 🤪`,
	success: `👁 Das Orakel wurde befragt und hat sich entschieden für: *%s*`,
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
	{"choose", nil},
	{"/choosex", nil},
	{" /choose", nil},
	{"/choose", expectedReply},
	{"/choose@bot", expectedReply},
	{"/choose foo* bar_ baz#", expectedReply},
	{"/choose@bot foo* bar_ baz#", expectedReply},
}

func provideMatcher() choose.Matcher {
	return choose.MakeMatcher()
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
			case strings.Contains((replies)[0].Text, "foo"):
				expectedReplies[0].Text = fmt.Sprintf(expected.success, `foo\*`)
			case strings.Contains((replies)[0].Text, "bar"):
				expectedReplies[0].Text = fmt.Sprintf(expected.success, `bar\_`)
			case strings.Contains((replies)[0].Text, "baz"):
				expectedReplies[0].Text = fmt.Sprintf(expected.success, `baz\#`)
			}

			assert.Equal(t, expectedReplies, replies, tt.in)
		}
	}
}
