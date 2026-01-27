package roll

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
)

const identifier = "roll"

const (
	defaultCount = 1
	defaultSides = 6
	maxCount     = 100
	maxSides     = 1000
)

const (
	statsTypeOverall = "overall"
	statsTypeLucky   = "lucky"
	statsTypeUnlucky = "unlucky"
	statsTypeUser    = "user"
)

const (
	topRollersLimit   = 5
	minThresholdValue = 0
	maxPartsLength    = 3
)

var (
	criticalHitMessages = []string{
		"Natural 20! The dice gods smile upon you!",
		"Perfect roll! Everything goes exactly right!",
		"Amazing! You couldn't have done better!",
		"Legendary success! Tales will be told of this roll!",
		"Jackpot! Maximum power achieved!",
	}

	criticalFailureMessages = []string{
		"Natural 1! Even your dice betray you!",
		"Catastrophic failure! Everything that could go wrong, did!",
		"Ouch! The dice gods laugh at your misfortune!",
		"Epic fail! You'll never live this down!",
		"Disaster! Your dice must hate you!",
	}
)

var pattern = regexp.MustCompile(`(?i)^/(roll)(@\w+)?($| )(.*)$`)

var help = []matcher.HelpStruct{
	{
		Command:     `roll`,
		Description: `Rolls dice in D&D notation with optional threshold and advantage mechanics.`,
		Usage:       `/roll [dice] [threshold] [kh]`,
		Example:     `/roll 2d20 15 kh`,
	},
	{
		Command:     `roll stats`,
		Description: `Shows roll statistics.`,
		Usage:       `/roll stats [user|lucky|unlucky]`,
		Example:     `/roll stats @username`,
	},
}

// Matcher implements the roll command matcher.
type Matcher struct {
	matcher.Matcher

	repo interfaces.RollRepoInterface
}

// MakeMatcher creates a new roll matcher instance.
func MakeMatcher(repo interfaces.RollRepoInterface) Matcher {
	return Matcher{
		Matcher: matcher.MakeMatcher(identifier, pattern, help),
		repo:    repo,
	}
}

// Process handles incoming messages matching the roll pattern.
func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	match := m.CommandMatch(messageIn)
	if match == nil {
		return nil, errors.New("message does not match")
	}

	// Extract and trim arguments
	args := strings.TrimSpace(match[3])

	// Check if it's a stats command
	if strings.HasPrefix(strings.ToLower(args), "stats") {
		return m.processStats(messageIn, args)
	}

	return m.processRoll(messageIn, args)
}

// processRoll handles the dice rolling command.
func (m Matcher) processRoll(messageIn telegramclient.WebhookMessageStruct, args string) ([]telegramclient.MessageStruct, error) {
	// Parse dice notation
	count, sides, threshold, keepHighest, err := parseDiceNotation(args)
	if err != nil {
		return m.makeReply("❌ "+telegramclient.EscapeMarkdown(err.Error()), messageIn.ID)
	}

	// Create and perform roll
	roll := NewDiceRoll(count, sides, threshold, keepHighest)
	roll.Roll()

	// Save to database
	if err := m.saveRoll(messageIn.From.ID, roll); err != nil {
		// Log error but don't fail the response
		// In production, you might want to log this properly
		_ = err
	}

	// Format response
	response := formatRollResponse(roll)

	return m.makeReply(response, messageIn.ID)
}

// processStats handles the stats command.
//
//nolint:cyclop // Complexity is acceptable for stats routing logic.
func (m Matcher) processStats(messageIn telegramclient.WebhookMessageStruct, args string) ([]telegramclient.MessageStruct, error) {
	// Remove "stats" prefix and trim
	args = strings.TrimSpace(strings.TrimPrefix(strings.ToLower(args), "stats"))

	var (
		stats      *interfaces.RollStatsStruct
		topRollers []interfaces.RollStatsStruct
	)

	statsType := statsTypeOverall

	var err error

	switch args {
	case "":
		// Overall stats
		stats, err = m.repo.GetOverallStats()
		if err != nil {
			return m.makeReply("❌ Error retrieving statistics\\.", messageIn.ID)
		}

		topRollers, _ = m.repo.GetTopRollers(topRollersLimit)
	case statsTypeLucky:
		// Luckiest roller
		stats, err = m.repo.GetLuckiestRoller()
		if err != nil || stats == nil || stats.TotalRolls == 0 {
			return m.makeReply("No lucky roller found \\(minimum 10 rolls required\\)\\.", messageIn.ID)
		}

		statsType = statsTypeLucky
	case statsTypeUnlucky:
		// Unluckiest roller
		stats, err = m.repo.GetUnluckiestRoller()
		if err != nil || stats == nil || stats.TotalRolls == 0 {
			return m.makeReply("No unlucky roller found \\(minimum 10 rolls required\\)\\.", messageIn.ID)
		}

		statsType = statsTypeUnlucky
	default:
		// User-specific stats
		// Remove @ if present
		username := strings.TrimPrefix(args, "@")

		// Look up user by username
		userID := m.findUserIDByUsername(username)
		if userID == 0 {
			return m.makeReply("❌ User not found: "+telegramclient.EscapeMarkdown(username), messageIn.ID)
		}

		stats, err = m.repo.GetUserStats(userID)
		if err != nil {
			return m.makeReply("❌ Error retrieving statistics\\.", messageIn.ID)
		}

		statsType = statsTypeUser
	}

	// Format response
	response := formatStatsResponse(stats, topRollers, statsType)

	return m.makeReply(response, messageIn.ID)
}

// saveRoll saves a roll to the database.
func (m Matcher) saveRoll(userID int64, roll *DiceRoll) error {
	// Convert results to JSON
	resultsJSON, err := json.Marshal(roll.GetResults())
	if err != nil {
		return err
	}

	// Create Roll struct
	dbRoll := &interfaces.Roll{
		UserID:          userID,
		DiceCount:       roll.GetCount(),
		DiceSides:       roll.GetSides(),
		Results:         string(resultsJSON),
		Total:           roll.Sum(),
		Threshold:       roll.GetThreshold(),
		KeepHighest:     roll.GetKeepHighest(),
		Success:         roll.IsSuccess(),
		CriticalHit:     roll.IsCriticalHit(),
		CriticalFailure: roll.IsCriticalFailure(),
	}

	return m.repo.SaveRoll(dbRoll)
}

// findUserIDByUsername looks up a user ID by username.
func (m Matcher) findUserIDByUsername(username string) int64 {
	userID, err := m.repo.GetUserIDByUsername(username)
	if err != nil {
		return 0
	}

	return userID
}

// parseDiceNotation parses dice notation string and returns count, sides, threshold, keepHighest, and error.
//
//nolint:funlen,cyclop,mnd // Dice notation parsing requires extensive validation
func parseDiceNotation(input string) (int, int, *int, bool, error) {
	input = strings.TrimSpace(input)

	// Default values
	if input == "" {
		return defaultCount, defaultSides, nil, false, nil
	}

	// Split by spaces
	parts := strings.Fields(input)

	// Parse dice notation (first part)
	diceRegex := regexp.MustCompile(`^(\d+)d(\d+)$`)

	matches := diceRegex.FindStringSubmatch(strings.ToLower(parts[0]))
	if matches == nil {
		return 0, 0, nil, false, errors.New("invalid dice notation! Use format: 2d20 or 2d20 15")
	}

	count, _ := strconv.Atoi(matches[1])
	sides, _ := strconv.Atoi(matches[2])

	// Validate dice values
	if count <= 0 {
		return 0, 0, nil, false, errors.New("dice count must be greater than 0")
	}

	if sides <= 1 {
		return 0, 0, nil, false, errors.New("dice sides must be greater than 1")
	}

	if count > maxCount {
		return 0, 0, nil, false, fmt.Errorf("maximum dice count is %d", maxCount)
	}

	if sides > maxSides {
		return 0, 0, nil, false, fmt.Errorf("maximum dice sides is %d", maxSides)
	}

	// Parse threshold (second part, optional)
	var threshold *int

	keepHighest := false

	if len(parts) >= 2 {
		t, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, nil, false, errors.New("invalid threshold value")
		}

		if t < minThresholdValue {
			return 0, 0, nil, false, errors.New("threshold must be non-negative")
		}

		threshold = &t
	}

	// Parse keep highest flag (third part, optional)
	if len(parts) >= 3 {
		if threshold == nil {
			return 0, 0, nil, false, errors.New("kh modifier requires a threshold")
		}

		if strings.ToLower(parts[2]) == "kh" {
			keepHighest = true
		} else {
			return 0, 0, nil, false, fmt.Errorf("unknown modifier: %s (use 'kh' for advantage)", parts[2])
		}
	}

	// Check for extra parameters
	if len(parts) > maxPartsLength {
		return 0, 0, nil, false, errors.New("too many parameters")
	}

	return count, sides, threshold, keepHighest, nil
}

// makeReply creates a reply message.
func (m Matcher) makeReply(message string, messageID int64) ([]telegramclient.MessageStruct, error) {
	return []telegramclient.MessageStruct{
		telegramclient.MarkdownReply(message, messageID),
	}, nil
}
