package ping

import (
	"testing"

	_ "github.com/br0-space/bot/testing_init"

	"github.com/br0-space/bot/internal/telegram/webhook"
	"github.com/stretchr/testify/assert"
)

func TestMakeMatcher(t *testing.T) {
	matcher := MakeMatcher()
	assert.NotNil(t, matcher)
	assert.IsType(t, Matcher{}, matcher)
}

func TestMatcher_Identifier(t *testing.T) {
	matcher := MakeMatcher()
	assert.Equal(t, matcher.Identifier(), "ping")
}

func TestMatcher_ProcessMessageNotMatching(t *testing.T) {
	message := webhook.MakeTestMessage(1, 2, "foobar")
	matcher := MakeMatcher()
	err := matcher.ProcessMessage(message)
	assert.Nil(t, err)
	assert.Len(t, matcher.telegram.SentMessages(), 0)
}

func TestMatcher_ProcessMessageMatching(t *testing.T) {
	message := webhook.MakeTestMessage(1, 2, "/ping")
	matcher := MakeMatcher()
	err := matcher.ProcessMessage(message)
	assert.Nil(t, err)
	assert.Len(t, matcher.telegram.SentMessages(), 1)
	sentMessage := matcher.telegram.SentMessages()[0]
	assert.NotNil(t, sentMessage)
	assert.Equal(t, sentMessage.ChatID, int64(1))
	assert.Equal(t, sentMessage.ReplyToMessageID, int64(2))
	assert.Equal(t, sentMessage.Text, "pong")
}