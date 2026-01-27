package interfaces

import "gorm.io/gorm"

// Roll represents a single dice roll in the database.
type Roll struct {
	gorm.Model `exhaustruct:"optional"`

	UserID          int64  `gorm:"not null;index"`
	DiceCount       int    `gorm:"not null"`
	DiceSides       int    `gorm:"not null"`
	Results         string `gorm:"type:text;not null"` // JSON array: "[5,12,18]"
	Total           int    `gorm:"not null"`
	Threshold       *int   `gorm:"index"`
	KeepHighest     bool   `gorm:"not null;default:false"`
	Success         *bool
	CriticalHit     bool `gorm:"not null;default:false;index"`
	CriticalFailure bool `gorm:"not null;default:false;index"`
}

// RollStatsStruct contains aggregated statistics for rolls.
type RollStatsStruct struct {
	UserID           int64
	Username         string
	TotalRolls       int
	TotalDice        int
	AverageRoll      float64
	HighestRoll      int
	LowestRoll       int
	CriticalHits     int
	CriticalFailures int
	SuccessRate      float64
}

// RollRepoInterface defines the repository interface for roll operations.
type RollRepoInterface interface {
	SaveRoll(roll *Roll) error
	GetOverallStats() (*RollStatsStruct, error)
	GetUserStats(userID int64) (*RollStatsStruct, error)
	GetUserIDByUsername(username string) (int64, error)
	GetLuckiestRoller() (*RollStatsStruct, error)
	GetUnluckiestRoller() (*RollStatsStruct, error)
	GetTopRollers(limit int) ([]RollStatsStruct, error)
}
