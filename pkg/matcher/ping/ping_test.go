package ping_test

import (
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matcher/ping"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedReply = []interfaces.TelegramMessageStruct{{
	ChatID:                0,
	Text:                  "pong",
	ParseMode:             "",
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
	{"ping", nil},
	{"/pings", nil},
	{" /ping", nil},
	{"/ping", expectedReply},
	{"/ping foo", expectedReply},
	{"/ping@bot", expectedReply},
	{"/ping@bot foo", expectedReply},
}

func provideMatcher() ping.Matcher {
	return ping.MakeMatcher()
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
			assert.Equal(t, tt.expectedReplies, replies, tt.in)
		}
	}
}
