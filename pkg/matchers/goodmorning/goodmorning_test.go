package goodmorning_test

import (
	"errors"
	"testing"
	"time"

	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matchers/goodmorning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Mocks

type mockState struct {
	lastPost *time.Time
}

func (m *mockState) ProcessMessage(_ telegramclient.WebhookMessageStruct) {}

func (m *mockState) GetLastPost(_ int64) *time.Time {
	return m.lastPost
}

type mockFortune struct{}

func (m *mockFortune) File() string       { return "test" }
func (m *mockFortune) ToMarkdown() string { return "fortune text" }

type mockFortuneService struct{}

func (m *mockFortuneService) GetList() []string                                      { return nil }
func (m *mockFortuneService) Exists(_ string) bool                                   { return false }
func (m *mockFortuneService) GetRandomFortune() (interfaces.FortuneInterface, error) { return &mockFortune{}, nil }
func (m *mockFortuneService) GetFortune(_ string) (interfaces.FortuneInterface, error) {
	return nil, errors.New("not implemented")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helpers

func isMorning() bool {
	h := time.Now().Hour()
	return h >= 6 && h <= 14
}

func newMatcher(lastPost *time.Time) goodmorning.Matcher {
	return goodmorning.MakeMatcher(&mockState{lastPost: lastPost}, &mockFortuneService{})
}

func newTestMessage() telegramclient.WebhookMessageStruct {
	return telegramclient.TestWebhookMessage("hello")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Tests

// TestProcess_NoReply_WhenLastPostIsNow is the core regression test for the processing-order bug.
// Before the fix, stateService.ProcessMessage() ran before matchers, updating lastPost to time.Now().
// The matcher then saw a sub-second delta and never fired. This test verifies that a just-now
// lastPost (the state ProcessMessage would leave behind) correctly suppresses the greeting.
func TestProcess_NoReply_WhenLastPostIsNow(t *testing.T) {
	t.Parallel()

	now := time.Now()
	replies, err := newMatcher(&now).Process(newTestMessage())

	require.NoError(t, err)
	assert.Nil(t, replies)
}

// TestProcess_NoReply_WhenLastPostIsRecent verifies that a post within the 6-hour window
// suppresses the greeting, regardless of time of day.
func TestProcess_NoReply_WhenLastPostIsRecent(t *testing.T) {
	t.Parallel()

	recent := time.Now().Add(-3 * time.Hour)
	replies, err := newMatcher(&recent).Process(newTestMessage())

	require.NoError(t, err)
	assert.Nil(t, replies)
}

// TestProcess_Replies_WhenNeverPosted verifies the matcher fires for a first-time user.
// Only meaningful during morning hours (06:00–14:00); skipped otherwise.
func TestProcess_Replies_WhenNeverPosted(t *testing.T) {
	t.Parallel()

	if !isMorning() {
		t.Skip("goodmorning matcher only fires between 06:00 and 14:00")
	}

	replies, err := newMatcher(nil).Process(newTestMessage())

	require.NoError(t, err)
	assert.NotNil(t, replies)
	assert.Len(t, replies, 1)
}

// TestProcess_Replies_WhenLastPostWasLongAgo verifies the matcher fires when the user's
// previous post was more than 6 hours ago.
// Only meaningful during morning hours (06:00–14:00); skipped otherwise.
func TestProcess_Replies_WhenLastPostWasLongAgo(t *testing.T) {
	t.Parallel()

	if !isMorning() {
		t.Skip("goodmorning matcher only fires between 06:00 and 14:00")
	}

	longAgo := time.Now().Add(-7 * time.Hour)
	replies, err := newMatcher(&longAgo).Process(newTestMessage())

	require.NoError(t, err)
	assert.NotNil(t, replies)
	assert.Len(t, replies, 1)
}
