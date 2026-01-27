package repo

import (
	"sync"

	"github.com/br0-space/bot/interfaces"
	"gorm.io/gorm"
)

var mutexRoll sync.Mutex

// RollRepo implements the RollRepoInterface for database operations.
type RollRepo struct {
	BaseRepo
}

// NewRollRepo creates a new RollRepo instance.
func NewRollRepo(tx *gorm.DB) *RollRepo {
	return &RollRepo{
		BaseRepo: NewBaseRepo(
			tx,
			&interfaces.Roll{},
		),
	}
}

// SaveRoll saves a roll to the database.
func (r RollRepo) SaveRoll(roll *interfaces.Roll) error {
	mutexRoll.Lock()
	defer mutexRoll.Unlock()

	return r.tx.Create(roll).Error
}

// GetOverallStats returns overall statistics for all rolls.
func (r RollRepo) GetOverallStats() (*interfaces.RollStatsStruct, error) {
	var stats interfaces.RollStatsStruct

	err := r.tx.Model(&interfaces.Roll{}).
		Select(`
			COUNT(*) as total_rolls,
			SUM(dice_count) as total_dice,
			AVG(total * 1.0 / dice_count) as average_roll,
			MAX(total) as highest_roll,
			MIN(total) as lowest_roll,
			SUM(CASE WHEN critical_hit THEN 1 ELSE 0 END) as critical_hits,
			SUM(CASE WHEN critical_failure THEN 1 ELSE 0 END) as critical_failures,
			AVG(CASE WHEN success IS NOT NULL AND success = true THEN 100.0 ELSE 0 END) as success_rate
		`).
		Where("deleted_at IS NULL").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetUserStats returns statistics for a specific user.
func (r RollRepo) GetUserStats(userID int64) (*interfaces.RollStatsStruct, error) {
	var stats interfaces.RollStatsStruct

	err := r.tx.Model(&interfaces.Roll{}).
		Select(`
			r.user_id,
			s.username,
			COUNT(*) as total_rolls,
			SUM(r.dice_count) as total_dice,
			AVG(r.total * 1.0 / r.dice_count) as average_roll,
			MAX(r.total) as highest_roll,
			MIN(r.total) as lowest_roll,
			SUM(CASE WHEN r.critical_hit THEN 1 ELSE 0 END) as critical_hits,
			SUM(CASE WHEN r.critical_failure THEN 1 ELSE 0 END) as critical_failures,
			AVG(CASE WHEN r.success IS NOT NULL AND r.success = true THEN 100.0 ELSE 0 END) as success_rate
		`).
		Table("rolls r").
		Joins("INNER JOIN stats s ON r.user_id = s.user_id").
		Where("r.user_id = ? AND r.deleted_at IS NULL", userID).
		Group("r.user_id, s.username").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetUserIDByUsername looks up a user ID by username (case-insensitive exact match).
func (r RollRepo) GetUserIDByUsername(username string) (int64, error) {
	var stats interfaces.Stats

	err := r.tx.Model(&interfaces.Stats{}).
		Where("LOWER(username) = LOWER(?)", username).
		First(&stats).Error
	if err != nil {
		return 0, err
	}

	return stats.UserID, nil
}

// GetLuckiestRoller returns the user with the highest average roll.
func (r RollRepo) GetLuckiestRoller() (*interfaces.RollStatsStruct, error) {
	var stats interfaces.RollStatsStruct

	err := r.tx.Model(&interfaces.Roll{}).
		Select(`
			r.user_id,
			s.username,
			COUNT(*) as total_rolls,
			SUM(r.dice_count) as total_dice,
			AVG(r.total * 1.0 / r.dice_count) as average_roll,
			MAX(r.total) as highest_roll,
			MIN(r.total) as lowest_roll,
			SUM(CASE WHEN r.critical_hit THEN 1 ELSE 0 END) as critical_hits,
			SUM(CASE WHEN r.critical_failure THEN 1 ELSE 0 END) as critical_failures
		`).
		Table("rolls r").
		Joins("INNER JOIN stats s ON r.user_id = s.user_id").
		Where("r.deleted_at IS NULL").
		Group("r.user_id, s.username").
		Having("COUNT(*) >= 10"). // Minimum rolls to be considered
		Order("average_roll DESC").
		Limit(1).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetUnluckiestRoller returns the user with the lowest average roll.
func (r RollRepo) GetUnluckiestRoller() (*interfaces.RollStatsStruct, error) {
	var stats interfaces.RollStatsStruct

	err := r.tx.Model(&interfaces.Roll{}).
		Select(`
			r.user_id,
			s.username,
			COUNT(*) as total_rolls,
			SUM(r.dice_count) as total_dice,
			AVG(r.total * 1.0 / r.dice_count) as average_roll,
			MAX(r.total) as highest_roll,
			MIN(r.total) as lowest_roll,
			SUM(CASE WHEN r.critical_hit THEN 1 ELSE 0 END) as critical_hits,
			SUM(CASE WHEN r.critical_failure THEN 1 ELSE 0 END) as critical_failures
		`).
		Table("rolls r").
		Joins("INNER JOIN stats s ON r.user_id = s.user_id").
		Where("r.deleted_at IS NULL").
		Group("r.user_id, s.username").
		Having("COUNT(*) >= 10"). // Minimum rolls to be considered
		Order("average_roll ASC").
		Limit(1).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetTopRollers returns the top N most active rollers.
func (r RollRepo) GetTopRollers(limit int) ([]interfaces.RollStatsStruct, error) {
	var stats []interfaces.RollStatsStruct

	err := r.tx.Model(&interfaces.Roll{}).
		Select(`
			r.user_id,
			s.username,
			COUNT(*) as total_rolls
		`).
		Table("rolls r").
		Joins("INNER JOIN stats s ON r.user_id = s.user_id").
		Where("r.deleted_at IS NULL").
		Group("r.user_id, s.username").
		Order("total_rolls DESC").
		Limit(limit).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats, nil
}
