package roll_test

import (
	"strings"
	"testing"

	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/matchers/roll"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Format Tests (using exported functions indirectly through Process)
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestFormatRollResponse_Basic(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll 2d20"))
	require.NoError(t, err)
	assert.Len(t, replies, 1)

	text := replies[0].Text

	// Should contain basic elements
	assert.Contains(t, text, "🎲")
	assert.Contains(t, text, "2d20")
	assert.Contains(t, text, "Results:")
	assert.Contains(t, text, "Total:")
}

func TestFormatRollResponse_WithAdvantage(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll 2d20 15 kh"))
	require.NoError(t, err)
	assert.Len(t, replies, 1)

	text := replies[0].Text

	// Should contain advantage indicator
	assert.Contains(t, text, "Advantage")
	assert.Contains(t, text, "Highest:")
	// Should not say "Total:" in advantage mode
	assert.NotContains(t, text, "Total:")
}

func TestFormatRollResponse_WithSuccess(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	matcher := roll.MakeMatcher(mockRepo)

	// Run multiple times to get both success and failure
	foundSuccess := false
	foundFailure := false

	for i := 0; i < 100 && (!foundSuccess || !foundFailure); i++ {
		replies, err := matcher.Process(newTestMessage("/roll 2d20 20"))
		require.NoError(t, err)
		assert.Len(t, replies, 1)

		text := replies[0].Text

		// Check for success
		if strings.Contains(text, "Success") {
			foundSuccess = true

			assert.Contains(t, text, "✅")
			assert.Contains(t, text, "Threshold: 20")
		}

		// Check for failure
		if strings.Contains(text, "Failure") {
			foundFailure = true

			assert.Contains(t, text, "❌")
			assert.Contains(t, text, "Threshold: 20")
		}
	}
}

func TestFormatStatsResponse_NoData(t *testing.T) {
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
	assert.Len(t, replies, 1)

	assert.Contains(t, replies[0].Text, "No roll")
}

func TestFormatStatsResponse_Overall(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetOverallStats").Return(&interfaces.RollStatsStruct{
		TotalRolls:       1234,
		TotalDice:        5678,
		AverageRoll:      10.5,
		HighestRoll:      60,
		LowestRoll:       2,
		CriticalHits:     15,
		CriticalFailures: 8,
		SuccessRate:      67.5,
	}, nil)
	mockRepo.On("GetTopRollers", 5).Return([]interfaces.RollStatsStruct{
		{Username: "user1", TotalRolls: 500},
		{Username: "user2", TotalRolls: 300},
		{Username: "user3", TotalRolls: 200},
	}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats"))
	require.NoError(t, err)
	assert.Len(t, replies, 1)

	text := replies[0].Text

	// Should contain all stats
	assert.Contains(t, text, "📊")
	assert.Contains(t, text, "1,234") // Formatted number with comma
	assert.Contains(t, text, "5,678")
	assert.Contains(t, text, "10.5")
	assert.Contains(t, text, "60")
	assert.Contains(t, text, "2")
	assert.Contains(t, text, "15")
	assert.Contains(t, text, "8")
	assert.Contains(t, text, "67.5")

	// Should contain top rollers
	assert.Contains(t, text, "Top Rollers")
	assert.Contains(t, text, "user1")
	assert.Contains(t, text, "user2")
	assert.Contains(t, text, "user3")
}

func TestFormatStatsResponse_User(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetUserIDByUsername", "testuser").Return(int64(123), nil)
	mockRepo.On("GetUserStats", int64(123)).Return(&interfaces.RollStatsStruct{
		Username:         "testuser",
		TotalRolls:       50,
		TotalDice:        125,
		AverageRoll:      11.2,
		HighestRoll:      40,
		LowestRoll:       3,
		CriticalHits:     2,
		CriticalFailures: 1,
		SuccessRate:      72.5,
	}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats testuser"))
	require.NoError(t, err)
	assert.Len(t, replies, 1)

	text := replies[0].Text

	// Should contain username in header
	assert.Contains(t, text, "testuser")
	assert.Contains(t, text, "50")
	assert.Contains(t, text, "125")
	assert.Contains(t, text, "11.2")
}

func TestFormatStatsResponse_Lucky(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetLuckiestRoller").Return(&interfaces.RollStatsStruct{
		Username:    "luckyuser",
		TotalRolls:  75,
		AverageRoll: 15.8,
	}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats lucky"))
	require.NoError(t, err)
	assert.Len(t, replies, 1)

	text := replies[0].Text

	assert.Contains(t, text, "🍀")
	assert.Contains(t, text, "Luckiest")
	assert.Contains(t, text, "luckyuser")
	assert.Contains(t, text, "15.8")
	assert.Contains(t, text, "dice gods")
}

func TestFormatStatsResponse_Unlucky(t *testing.T) {
	t.Parallel()

	mockRepo := new(MockRollRepo)
	mockRepo.On("SaveRoll", mock.Anything).Return(nil)
	mockRepo.On("GetUnluckiestRoller").Return(&interfaces.RollStatsStruct{
		Username:    "unluckyuser",
		TotalRolls:  80,
		AverageRoll: 4.2,
	}, nil)
	matcher := roll.MakeMatcher(mockRepo)

	replies, err := matcher.Process(newTestMessage("/roll stats unlucky"))
	require.NoError(t, err)
	assert.Len(t, replies, 1)

	text := replies[0].Text

	assert.Contains(t, text, "💀")
	assert.Contains(t, text, "Unluckiest")
	assert.Contains(t, text, "unluckyuser")
	assert.Contains(t, text, "4.2")
	assert.Contains(t, text, "different dice")
}

func TestFormatNumber_Variations(t *testing.T) {
	t.Parallel()

	// Test through stats formatting
	tests := []struct {
		name     string
		rolls    int
		expected string
	}{
		{"small", 42, "42"},
		{"hundreds", 567, "567"},
		{"thousands", 1234, "1,234"},
		{"ten thousands", 12345, "12,345"},
		{"hundred thousands", 123456, "123,456"},
		{"millions", 1234567, "1,234,567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(MockRollRepo)
			mockRepo.On("SaveRoll", mock.Anything).Return(nil)
			mockRepo.On("GetOverallStats").Return(&interfaces.RollStatsStruct{
				TotalRolls: tt.rolls,
				TotalDice:  tt.rolls * 2,
			}, nil)
			mockRepo.On("GetTopRollers", 5).Return([]interfaces.RollStatsStruct{}, nil)
			matcher := roll.MakeMatcher(mockRepo)

			replies, err := matcher.Process(newTestMessage("/roll stats"))
			require.NoError(t, err)
			assert.Len(t, replies, 1)

			assert.Contains(t, replies[0].Text, tt.expected)
		})
	}
}
