package roll

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
)

const (
	percentMultiplier  = 100.0
	numberGroupSize    = 3
	minNumberGroupSize = 3
)

// formatRollResponse creates a formatted response message for a dice roll.
func formatRollResponse(roll *DiceRoll) string {
	var parts []string

	// Header
	if roll.GetKeepHighest() {
		parts = append(parts, fmt.Sprintf("🎲 *Roll: %dd%d \\(Advantage\\)*", roll.GetCount(), roll.GetSides()))
	} else {
		parts = append(parts, fmt.Sprintf("🎲 *Roll: %dd%d*", roll.GetCount(), roll.GetSides()))
	}

	// Results
	results := roll.GetResults()

	resultStrs := make([]string, len(results))
	for i, r := range results {
		resultStrs[i] = strconv.Itoa(r)
	}

	parts = append(parts, "Results: "+strings.Join(resultStrs, ", "))

	// Total or Highest
	if roll.GetKeepHighest() {
		parts = append(parts, fmt.Sprintf("Highest: *%d*", roll.GetHighest()))
	} else {
		parts = append(parts, fmt.Sprintf("Total: *%d*", roll.Sum()))
	}

	// Success/Failure
	if success := roll.IsSuccess(); success != nil {
		if *success {
			parts = append(parts, fmt.Sprintf("✅ Success\\! \\(Threshold: %d\\)", *roll.GetThreshold()))
		} else {
			parts = append(parts, fmt.Sprintf("❌ Failure\\! \\(Threshold: %d\\)", *roll.GetThreshold()))
		}
	}

	// Critical hit/failure messages
	if roll.IsCriticalHit() {
		parts = append(parts, "✨ *CRITICAL HIT\\!* ✨")
		parts = append(parts, "🎉 "+telegramclient.EscapeMarkdown(getRandomCriticalHit()))
	} else if roll.IsCriticalFailure() {
		parts = append(parts, "💀 *CRITICAL FAILURE\\!* 💀")
		parts = append(parts, "😂 "+telegramclient.EscapeMarkdown(getRandomCriticalFailure()))
	}

	return strings.Join(parts, "\n")
}

// getRandomCriticalHit returns a random critical hit message.
func getRandomCriticalHit() string {
	if len(criticalHitMessages) == 0 {
		return "Critical hit!"
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(criticalHitMessages))))

	return criticalHitMessages[int(n.Int64())]
}

// getRandomCriticalFailure returns a random critical failure message.
func getRandomCriticalFailure() string {
	if len(criticalFailureMessages) == 0 {
		return "Critical failure!"
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(criticalFailureMessages))))

	return criticalFailureMessages[int(n.Int64())]
}

// formatStatsResponse creates a formatted response message for roll statistics.
//
//nolint:cyclop,funlen
func formatStatsResponse(stats *interfaces.RollStatsStruct, topRollers []interfaces.RollStatsStruct, statsType string) string {
	if stats == nil || stats.TotalRolls == 0 {
		if statsType == statsTypeUser {
			return "No roll data found for this user\\."
		}

		if statsType == statsTypeLucky {
			return "No lucky roller found \\(minimum 10 rolls required\\)\\."
		}

		if statsType == statsTypeUnlucky {
			return "No unlucky roller found \\(minimum 10 rolls required\\)\\."
		}

		return "No roll data available yet\\. Start rolling\\!"
	}

	var parts []string

	// Header
	switch statsType {
	case statsTypeUser:
		parts = append(parts, fmt.Sprintf("📊 *Roll Statistics for @%s*", telegramclient.EscapeMarkdown(stats.Username)))
	case statsTypeLucky:
		parts = append(parts, "🍀 *Luckiest Roller*")
		parts = append(parts, "")
		parts = append(parts, "@"+telegramclient.EscapeMarkdown(stats.Username))
	case statsTypeUnlucky:
		parts = append(parts, "💀 *Unluckiest Roller*")
		parts = append(parts, "")
		parts = append(parts, "@"+telegramclient.EscapeMarkdown(stats.Username))
	default:
		parts = append(parts, "📊 *Roll Statistics*")
	}

	parts = append(parts, "")

	// Main stats
	parts = append(parts, "Total Rolls: "+formatNumber(stats.TotalRolls))
	parts = append(parts, "Total Dice: "+formatNumber(stats.TotalDice))
	parts = append(parts, fmt.Sprintf("Average per Die: %.1f", stats.AverageRoll))

	if stats.HighestRoll > 0 || stats.LowestRoll > 0 {
		parts = append(parts, "")
		parts = append(parts, fmt.Sprintf("Highest Roll: %d", stats.HighestRoll))
		parts = append(parts, fmt.Sprintf("Lowest Roll: %d", stats.LowestRoll))
	}

	if stats.CriticalHits > 0 || stats.CriticalFailures > 0 {
		parts = append(parts, "")

		critHitRate := 0.0
		if stats.TotalRolls > 0 {
			critHitRate = float64(stats.CriticalHits) * percentMultiplier / float64(stats.TotalRolls)
		}

		parts = append(parts, fmt.Sprintf("Critical Hits: %d \\(%.1f%%\\)", stats.CriticalHits, critHitRate))

		critFailRate := 0.0
		if stats.TotalRolls > 0 {
			critFailRate = float64(stats.CriticalFailures) * percentMultiplier / float64(stats.TotalRolls)
		}

		parts = append(parts, fmt.Sprintf("Critical Failures: %d \\(%.1f%%\\)", stats.CriticalFailures, critFailRate))
	}

	if stats.SuccessRate > 0 {
		parts = append(parts, fmt.Sprintf("Success Rate: %.1f%%", stats.SuccessRate))
	}

	// Top rollers for overall stats
	if statsType == statsTypeOverall && len(topRollers) > 0 {
		parts = append(parts, "")

		parts = append(parts, "🏆 *Top Rollers:*")
		for i, roller := range topRollers {
			parts = append(parts, fmt.Sprintf("%d\\. @%s \\- %s rolls",
				i+1,
				telegramclient.EscapeMarkdown(roller.Username),
				formatNumber(roller.TotalRolls)))
		}
	}

	// Footer message for lucky/unlucky
	switch statsType {
	case statsTypeLucky:
		parts = append(parts, "")
		parts = append(parts, "The dice gods clearly favor you\\! 🎲✨")
	case statsTypeUnlucky:
		parts = append(parts, "")
		parts = append(parts, "Maybe try different dice? 🎲💀")
	}

	return strings.Join(parts, "\n")
}

// formatNumber formats a number with thousands separator.
func formatNumber(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= minNumberGroupSize {
		return s
	}

	// Add comma separators
	var result []string

	for i := len(s); i > 0; i -= numberGroupSize {
		start := max(i-numberGroupSize, 0)

		result = append([]string{s[start:i]}, result...)
	}

	return strings.Join(result, ",")
}
