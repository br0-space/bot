package ping_test

import (
	"testing"

	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/pkg/matchers/ping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedReply = []telegramclient.MessageStruct{{
	ChatID:                0,
	ReplyToMessageID:      123,
	Text:                  "pong",
	Photo:                 "",
	Caption:               "",
	ParseMode:             "",
	DisableWebPagePreview: false,
	DisableNotification:   false,
}}

var tests = []struct {
	in              string
	expectedReplies []telegramclient.MessageStruct
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
			assert.Equal(t, tt.expectedReplies, replies, tt.in)
		}
	}
}
