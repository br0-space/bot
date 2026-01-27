package roll_test

import (
	"strings"
	"testing"

	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matchers/roll"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRollRepo is a mock implementation of RollRepoInterface.
type MockRollRepo struct {
	mock.Mock
}

func (m *MockRollRepo) SaveRoll(r *interfaces.Roll) error {
	args := m.Called(r)

	return args.Error(0)
}

func (m *MockRollRepo) GetOverallStats() (*interfaces.RollStatsStruct, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	stats, ok := args.Get(0).(*interfaces.RollStatsStruct)
	if !ok {
		return nil, args.Error(1)
	}

	return stats, args.Error(1)
}

func (m *MockRollRepo) GetUserStats(userID int64) (*interfaces.RollStatsStruct, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	stats, ok := args.Get(0).(*interfaces.RollStatsStruct)

	if !ok {
		return nil, args.Error(1)
	}

	return stats, args.Error(1)
}

func (m *MockRollRepo) GetLuckiestRoller() (*interfaces.RollStatsStruct, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	stats, ok := args.Get(0).(*interfaces.RollStatsStruct)
	if !ok {
		return nil, args.Error(1)
	}

	return stats, args.Error(1)
}

func (m *MockRollRepo) GetUnluckiestRoller() (*interfaces.RollStatsStruct, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	stats, ok := args.Get(0).(*interfaces.RollStatsStruct)
	if !ok {
		return nil, args.Error(1)
	}

	return stats, args.Error(1)
}

func (m *MockRollRepo) GetTopRollers(limit int) ([]interfaces.RollStatsStruct, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	rollers, ok := args.Get(0).([]interfaces.RollStatsStruct)
	if !ok {
		return nil, args.Error(1)
	}

	return rollers, args.Error(1)
}

func (m *MockRollRepo) GetUserIDByUsername(username string) (int64, error) {
	args := m.Called(username)

	userID, ok := args.Get(0).(int64)
	if !ok {
		return 0, args.Error(1)
	}

	return userID, args.Error(1)
}

func provideMatcher() roll.Matcher {
	mockRepo := new(MockRollRepo)
	// Setup default mock behavior - save rolls successfully
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)

	return roll.MakeMatcher(mockRepo)
}

func newTestMessage(text string) telegramclient.WebhookMessageStruct {
	return telegramclient.TestWebhookMessage(text)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// DoesMatch Tests
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var doesMatchTests = []struct {
	in            string
	expectedMatch bool
}{
	// Non-matches
	{"", false},
	{"foobar", false},
	{"roll", false},
	{"/rollx", false},
	{" /roll", false},
	{"hello /roll", false},

	// Basic matches
	{"/roll", true},
	{"/roll ", true},
	{"/roll@bot", true},
	{"/roll@bot ", true},
	{"/ROLL", true},
	{"/Roll", true},

	// Roll with arguments
	{"/roll 2d20", true},
	{"/roll 1d6", true},
	{"/roll 3d8", true},
	{"/roll 2d20 15", true},
	{"/roll 3d6 10", true},
	{"/roll 2d20 15 kh", true},
	{"/roll 4d6 12 KH", true},
	{"/roll@bot 2d20", true},
	{"/roll@bot 2d20 15", true},
	{"/roll@bot 2d20 15 kh", true},

	// Stats commands
	{"/roll stats", true},
	{"/roll stats ", true},
	{"/roll STATS", true},
	{"/roll Stats", true},
	{"/roll stats lucky", true},
	{"/roll stats unlucky", true},
	{"/roll stats @username", true},
	{"/roll stats username", true},
	{"/roll@bot stats", true},
	{"/roll@bot stats lucky", true},
	{"/roll@bot stats @user", true},
}

func TestMatcher_DoesMatch(t *testing.T) {
	t.Parallel()

	matcher := provideMatcher()
	for _, tt := range doesMatchTests {
		doesMatch := matcher.DoesMatch(newTestMessage(tt.in))
		assert.Equal(t, tt.expectedMatch, doesMatch, "Input: %s", tt.in)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Process Tests - Basic Roll Commands
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestMatcher_Process_DefaultRoll(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)

	// Should contain dice emoji and result
	assert.Contains(t, replies[0].Text, "🎲")
	// Default is 1d6
	assert.Contains(t, replies[0].Text, "Results:")
	assert.Regexp(t, `Results: \d`, replies[0].Text)
}

func TestMatcher_Process_CustomDiceRoll(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll 2d20"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)

	// Should contain dice emoji and two results
	assert.Contains(t, replies[0].Text, "🎲")
	assert.Contains(t, replies[0].Text, "2d20")
	// Should have total
	assert.Contains(t, replies[0].Text, "Total:")
}

func TestMatcher_Process_WithThreshold(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll 2d20 15"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)

	// Should contain threshold in output
	assert.Contains(t, replies[0].Text, "15")
	// Should indicate success or failure
	text := replies[0].Text
	hasSuccess := strings.Contains(text, "Success")
	hasFailure := strings.Contains(text, "Failure")
	assert.True(t, hasSuccess || hasFailure, "Response should contain Success or Failure")
}

func TestMatcher_Process_WithKeepHighest(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll 2d20 15 kh"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)

	// Should contain advantage indicator
	assert.Contains(t, replies[0].Text, "Advantage")
	// Should indicate success or failure
	text := replies[0].Text
	hasSuccess := strings.Contains(text, "Success")
	hasFailure := strings.Contains(text, "Failure")
	assert.True(t, hasSuccess || hasFailure, "Response should contain Success or Failure")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Process Tests - Invalid Input
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestMatcher_Process_InvalidDiceNotation(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	invalidInputs := []string{
		"/roll d20",      // No count
		"/roll 2d",       // No sides
		"/roll abc",      // Invalid format
		"/roll 2d20d6",   // Too many d's
		"/roll 0d6",      // Zero count
		"/roll 2d0",      // Zero sides
		"/roll 2d1",      // Only 1 side
		"/roll -2d6",     // Negative count
		"/roll 2d-6",     // Negative sides
		"/roll 1000d20",  // Too many dice
		"/roll 2d10000",  // Too many sides
		"/roll 2d20 -5",  // Negative threshold
		"/roll 2d20 abc", // Invalid threshold
	}

	for _, input := range invalidInputs {
		replies, err := matcher.Process(newTestMessage(input))
		require.NoError(t, err, "Input: %s", input)
		require.NotNil(t, replies, "Input: %s", input)
		assert.Contains(t, replies[0].Text, "❌", "Input: %s should return error message", input)
	}
}

func TestMatcher_Process_KeepHighestWithoutThreshold(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll 2d20 kh"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Contains(t, replies[0].Text, "❌")
	assert.Contains(t, replies[0].Text, "invalid threshold")
}

func TestMatcher_Process_InvalidModifier(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll 2d20 15 invalid"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Contains(t, replies[0].Text, "❌")
	assert.Contains(t, replies[0].Text, "unknown modifier")
}

func TestMatcher_Process_TooManyParameters(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll 2d20 15 kh extra"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Contains(t, replies[0].Text, "❌")
	assert.Contains(t, replies[0].Text, "too many parameters")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Process Tests - Stats Commands
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestMatcher_Process_StatsOverall_NoData(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetOverallStats").Return(&interfaces.RollStatsStruct{
		TotalRolls: 0,
	}, nil)
	mockRepo.On("GetTopRollers", 5).Return([]interfaces.RollStatsStruct{}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)

	// Should indicate no data
	assert.Contains(t, replies[0].Text, "No roll data available yet")
}

func TestMatcher_Process_StatsOverall_WithData(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetOverallStats").Return(&interfaces.RollStatsStruct{
		TotalRolls:       100,
		TotalDice:        250,
		AverageRoll:      10.5,
		HighestRoll:      40,
		LowestRoll:       2,
		CriticalHits:     5,
		CriticalFailures: 3,
		SuccessRate:      65.0,
	}, nil)
	mockRepo.On("GetTopRollers", 5).Return([]interfaces.RollStatsStruct{
		{Username: "user1", TotalRolls: 50},
		{Username: "user2", TotalRolls: 30},
	}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)

	// Should contain statistics
	assert.Contains(t, replies[0].Text, "Roll Statistics")
	assert.Contains(t, replies[0].Text, "100")  // Total rolls
	assert.Contains(t, replies[0].Text, "10.5") // Average roll
}

func TestMatcher_Process_StatsLucky(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetLuckiestRoller").Return(&interfaces.RollStatsStruct{
		Username:    "luckyuser",
		TotalRolls:  50,
		AverageRoll: 15.2,
	}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats lucky"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)

	assert.Contains(t, replies[0].Text, "Luckiest")
	assert.Contains(t, replies[0].Text, "luckyuser")
	assert.Contains(t, replies[0].Text, "15.2")
}

func TestMatcher_Process_StatsUnlucky(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetUnluckiestRoller").Return(&interfaces.RollStatsStruct{
		Username:    "unluckyuser",
		TotalRolls:  50,
		AverageRoll: 5.2,
	}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats unlucky"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)

	assert.Contains(t, replies[0].Text, "Unluckiest")
	assert.Contains(t, replies[0].Text, "unluckyuser")
	assert.Contains(t, replies[0].Text, "5.2")
}

func TestMatcher_Process_StatsUser(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetUserIDByUsername", "testuser").Return(int64(123), nil)
	mockRepo.On("GetUserStats", int64(123)).Return(&interfaces.RollStatsStruct{
		Username:    "testuser",
		TotalRolls:  25,
		AverageRoll: 11.5,
	}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	// Test with @ prefix
	replies, err := matcher.Process(newTestMessage("/roll stats @testuser"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)
	assert.Contains(t, replies[0].Text, "testuser")
	assert.Contains(t, replies[0].Text, "25")

	// Test without @ prefix
	replies, err = matcher.Process(newTestMessage("/roll stats testuser"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)
	assert.Contains(t, replies[0].Text, "testuser")
}

func TestMatcher_Process_StatsUser_NotFound(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetUserIDByUsername", "nonexistent").Return(int64(0), nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats nonexistent"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)
	assert.Contains(t, replies[0].Text, "❌")
	assert.Contains(t, replies[0].Text, "User not found")
}

func TestMatcher_Process_StatsLucky_NoData(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetLuckiestRoller").Return(&interfaces.RollStatsStruct{
		TotalRolls: 0,
	}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats lucky"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)
	assert.Contains(t, replies[0].Text, "No lucky roller found")
}

func TestMatcher_Process_StatsUnlucky_NoData(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetUnluckiestRoller").Return(&interfaces.RollStatsStruct{
		TotalRolls: 0,
	}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats unlucky"))
	require.NoError(t, err)
	require.NotNil(t, replies)
	assert.Len(t, replies, 1)
	assert.Contains(t, replies[0].Text, "No unlucky roller found")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Non-matching messages
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestMatcher_Process_NonMatch(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("just some text"))
	require.Error(t, err)
	assert.Nil(t, replies)
	assert.Contains(t, err.Error(), "does not match")
}
