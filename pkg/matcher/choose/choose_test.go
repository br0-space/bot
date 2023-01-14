package choose_test

import (
	"fmt"
	"github.com/br0-space/bot/container"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/choose"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var expected = struct {
	insult  string
	success string
}{
	insult:  `Ob du behindert bist hab ich gefragt?\! 🤪`,
	success: `👁 Das Orakel wurde befragt und hat sich entschieden für: *%s*`,
}

var expectedReply = []interfaces.TelegramMessageStruct{{
	ChatID:                0,
	Text:                  "",
	ParseMode:             "MarkdownV2",
	DisableWebPagePreview: false,
	DisableNotification:   false,
	ReplyToMessageID:      123,
}}

var tests = []struct {
	in              string
	expectedReplies []interfaces.TelegramMessageStruct
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
	return choose.MakeMatcher(
		container.ProvideLogger(),
	)
}

func newTestMessage(text string) interfaces.TelegramWebhookMessageStruct {
	return interfaces.NewTestTelegramWebhookMessage(text)
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
			assert.Error(t, err, tt.in)
			assert.Nil(t, replies, tt.in)
		} else {
			assert.NoError(t, err, tt.in)
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
