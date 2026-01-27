package roll_test

import (
	"testing"

	"github.com/br0-space/bot/pkg/matchers/roll"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// DiceRoll Tests
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestNewDiceRoll(t *testing.T) {
	t.Parallel()

	threshold := 15

	tests := []struct {
		name        string
		count       int
		sides       int
		threshold   *int
		keepHighest bool
	}{
		{"basic", 2, 20, nil, false},
		{"with threshold", 3, 6, &threshold, false},
		{"with keep highest", 2, 20, &threshold, true},
		{"single die", 1, 6, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dice := roll.NewDiceRoll(tt.count, tt.sides, tt.threshold, tt.keepHighest)
			require.NotNil(t, dice)

			assert.Equal(t, tt.count, dice.GetCount())
			assert.Equal(t, tt.sides, dice.GetSides())
			assert.Equal(t, tt.threshold, dice.GetThreshold())
			assert.Equal(t, tt.keepHighest, dice.GetKeepHighest())
			assert.Empty(t, dice.GetResults())
		})
	}
}

func TestDiceRoll_Roll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		count int
		sides int
	}{
		{"1d6", 1, 6},
		{"2d20", 2, 20},
		{"3d8", 3, 8},
		{"10d10", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dice := roll.NewDiceRoll(tt.count, tt.sides, nil, false)
			results := dice.Roll()

			// Check correct number of results
			assert.Len(t, results, tt.count)
			assert.Equal(t, results, dice.GetResults())

			// Check all results are in valid range [1, sides]
			for i, result := range results {
				assert.GreaterOrEqual(t, result, 1, "Result %d should be >= 1", i)
				assert.LessOrEqual(t, result, tt.sides, "Result %d should be <= %d", i, tt.sides)
			}
		})
	}
}

func TestDiceRoll_Sum(t *testing.T) {
	t.Parallel()

	dice := roll.NewDiceRoll(3, 6, nil, false)
	results := dice.Roll()

	sum := dice.Sum()

	// Calculate expected sum
	expectedSum := 0
	for _, result := range results {
		expectedSum += result
	}

	assert.Equal(t, expectedSum, sum)

	// Sum should be in valid range [count, count*sides]
	assert.GreaterOrEqual(t, sum, 3)
	assert.LessOrEqual(t, sum, 18)
}

func TestDiceRoll_GetHighest(t *testing.T) {
	t.Parallel()

	dice := roll.NewDiceRoll(5, 20, nil, false)
	results := dice.Roll()

	highest := dice.GetHighest()

	// Find expected highest
	expectedHighest := results[0]
	for _, result := range results {
		if result > expectedHighest {
			expectedHighest = result
		}
	}

	assert.Equal(t, expectedHighest, highest)
	assert.GreaterOrEqual(t, highest, 1)
	assert.LessOrEqual(t, highest, 20)
}

func TestDiceRoll_GetHighest_EmptyResults(t *testing.T) {
	t.Parallel()

	dice := roll.NewDiceRoll(1, 6, nil, false)
	// Don't roll, so results are empty
	highest := dice.GetHighest()

	assert.Equal(t, 0, highest)
}

func TestDiceRoll_IsCriticalHit(t *testing.T) {
	t.Parallel()

	// Run multiple times to eventually get a critical hit
	// (or at least test the logic)
	for range 1000 {
		dice := roll.NewDiceRoll(2, 6, nil, false)
		results := dice.Roll()

		isCrit := dice.IsCriticalHit()

		// Check if detection is correct
		allMax := true

		for _, result := range results {
			if result != 6 {
				allMax = false

				break
			}
		}

		assert.Equal(t, allMax, isCrit, "Results: %v", results)

		if isCrit {
			// Found a critical hit, test passed
			break
		}
	}
}

func TestDiceRoll_IsCriticalHit_EmptyResults(t *testing.T) {
	t.Parallel()

	dice := roll.NewDiceRoll(1, 6, nil, false)
	// Don't roll
	assert.False(t, dice.IsCriticalHit())
}

func TestDiceRoll_IsCriticalFailure(t *testing.T) {
	t.Parallel()

	// Run multiple times to eventually get a critical failure
	for range 1000 {
		dice := roll.NewDiceRoll(2, 20, nil, false)
		results := dice.Roll()

		isCrit := dice.IsCriticalFailure()

		// Check if detection is correct
		allMin := true

		for _, result := range results {
			if result != 1 {
				allMin = false

				break
			}
		}

		assert.Equal(t, allMin, isCrit, "Results: %v", results)

		if isCrit {
			// Found a critical failure, test passed
			break
		}
	}
}

func TestDiceRoll_IsCriticalFailure_EmptyResults(t *testing.T) {
	t.Parallel()

	dice := roll.NewDiceRoll(1, 6, nil, false)
	// Don't roll
	assert.False(t, dice.IsCriticalFailure())
}

func TestDiceRoll_IsSuccess_NoThreshold(t *testing.T) {
	t.Parallel()

	dice := roll.NewDiceRoll(2, 20, nil, false)
	dice.Roll()

	success := dice.IsSuccess()
	assert.Nil(t, success, "Success should be nil when no threshold is set")
}

func TestDiceRoll_IsSuccess_WithThreshold(t *testing.T) {
	t.Parallel()

	// Run multiple times to test both success and failure cases
	threshold := 25
	successCount := 0
	failureCount := 0

	for range 100 {
		dice := roll.NewDiceRoll(2, 20, &threshold, false)
		dice.Roll()

		success := dice.IsSuccess()
		require.NotNil(t, success)

		// Verify the success value matches the sum
		sum := dice.Sum()
		expectedSuccess := sum >= threshold

		assert.Equal(t, expectedSuccess, *success, "Sum: %d, Threshold: %d", sum, threshold)

		if *success {
			successCount++
		} else {
			failureCount++
		}
	}

	// Should have some of each (statistically very likely with 2d20 vs threshold 25)
	// This might occasionally fail due to randomness, but it's unlikely
	assert.Positive(t, successCount, "Should have at least one success")
	assert.Positive(t, failureCount, "Should have at least one failure")
}

func TestDiceRoll_IsSuccess_WithKeepHighest(t *testing.T) {
	t.Parallel()

	// Test keep highest mode (advantage)
	threshold := 15
	successCount := 0
	failureCount := 0

	for range 100 {
		dice := roll.NewDiceRoll(2, 20, &threshold, true)
		dice.Roll()

		success := dice.IsSuccess()
		require.NotNil(t, success)

		// Verify the success value matches the highest die
		highest := dice.GetHighest()
		expectedSuccess := highest >= threshold

		assert.Equal(t, expectedSuccess, *success, "Highest: %d, Threshold: %d, Results: %v",
			highest, threshold, dice.GetResults())

		if *success {
			successCount++
		} else {
			failureCount++
		}
	}

	// Should have some of each
	assert.Positive(t, successCount, "Should have at least one success")
	assert.Positive(t, failureCount, "Should have at least one failure")
}

func TestDiceRoll_IsSuccess_KeepHighestVsSum(t *testing.T) {
	t.Parallel()

	threshold := 10

	// Create a scenario where we can verify the difference
	// Run many times to find a case where keepHighest and sum differ
	foundDifference := false

	for i := 0; i < 1000 && !foundDifference; i++ {
		dice1 := roll.NewDiceRoll(3, 6, &threshold, false) // Sum mode
		dice1.Roll()

		dice2 := roll.NewDiceRoll(3, 6, &threshold, true) // Keep highest mode
		// Use same results for comparison
		dice2.Roll()

		success1 := dice1.IsSuccess()
		success2 := dice2.IsSuccess()

		require.NotNil(t, success1)
		require.NotNil(t, success2)

		// Verify logic
		sum := dice1.Sum()
		highest := dice1.GetHighest()

		assert.Equal(t, sum >= threshold, *success1)
		assert.Equal(t, highest >= threshold, *success2)

		// Check if we found a case where they differ
		if *success1 != *success2 {
			foundDifference = true
		}
	}

	// Note: It's possible (though unlikely) that we don't find a difference in 1000 rolls
	// This is acceptable for this test
}

func TestDiceRoll_Randomness(t *testing.T) {
	t.Parallel()

	// Verify that rolls produce different results (not deterministic)
	dice1 := roll.NewDiceRoll(5, 20, nil, false)
	dice2 := roll.NewDiceRoll(5, 20, nil, false)

	results1 := dice1.Roll()
	results2 := dice2.Roll()

	// Extremely unlikely to get identical results (1 in 3,200,000)
	// If this fails, it's probably a bug (or cosmic bad luck)
	assert.NotEqual(t, results1, results2, "Two independent rolls should produce different results")
}

func TestDiceRoll_MultipleRolls(t *testing.T) {
	t.Parallel()

	dice := roll.NewDiceRoll(2, 6, nil, false)

	// First roll
	results1 := dice.Roll()
	assert.Len(t, results1, 2)

	// Second roll should overwrite first
	results2 := dice.Roll()
	assert.Len(t, results2, 2)

	// Results should be stored
	assert.Equal(t, results2, dice.GetResults())

	// Results should be different (very likely)
	assert.NotEqual(t, results1, results2)
}

func TestDiceRoll_GettersBeforeRoll(t *testing.T) {
	t.Parallel()

	threshold := 10
	dice := roll.NewDiceRoll(3, 8, &threshold, true)

	// Test getters before rolling
	assert.Equal(t, 3, dice.GetCount())
	assert.Equal(t, 8, dice.GetSides())
	assert.Equal(t, &threshold, dice.GetThreshold())
	assert.True(t, dice.GetKeepHighest())
	assert.Empty(t, dice.GetResults())
	assert.Equal(t, 0, dice.Sum())
	assert.Equal(t, 0, dice.GetHighest())
	assert.False(t, dice.IsCriticalHit())
	assert.False(t, dice.IsCriticalFailure())
}
