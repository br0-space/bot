package matcher

import (
	"errors"
	"testing"

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
	assert.Equal(t, matcher.Identifier(), "matcher")
}

func TestMatcher_HandleError(t *testing.T) {
	message := webhook.MakeTestMessage(1, 2, "foobar")
	matcher := MakeMatcher()
	matcher.HandleError(message, errors.New("test"))
	assert.Len(t, matcher.telegram.SentMessages(), 1)
	assert.Equal(t, matcher.telegram.SentMessages()[0].ChatID, int64(1))
	assert.Equal(t, matcher.telegram.SentMessages()[0].ReplyToMessageID, int64(2))
	assert.Equal(t, matcher.telegram.SentMessages()[0].Text, "⚠️ matcher test")
}