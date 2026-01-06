package repo

import (
	"sync"
	"time"

	"github.com/br0-space/bot/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var mutexStats sync.Mutex

type UserStatsRepo struct {
	BaseRepo
}

func NewUserStatsRepo(tx *gorm.DB) *UserStatsRepo {
	return &UserStatsRepo{
		BaseRepo: NewBaseRepo(
			tx,
			&interfaces.Stats{
				UserID:   0,
				Username: "",
				Posts:    0,
				LastPost: time.Time{},
			},
		),
	}
}

func (r UserStatsRepo) UpdateStats(userID int64, username string) error {
	mutexStats.Lock()
	defer mutexStats.Unlock()

	return r.tx.Clauses(clause.OnConflict{ //nolint:exhaustruct
		Columns: []clause.Column{{Name: "user_id"}}, //nolint:exhaustruct
		DoUpdates: clause.Assignments(map[string]any{
			"username":  username,
			"posts":     gorm.Expr("stats.posts + 1"),
			"last_post": time.Now(),
		}),
	}).Create(&interfaces.Stats{
		UserID:   userID,
		Username: username,
		Posts:    1,
		LastPost: time.Now(),
	}).Error
}

func (r UserStatsRepo) GetKnownUsers() ([]interfaces.StatsUserStruct, error) {
	var records []interfaces.Stats

	r.tx.
		Where("user_id != 0").
		Order("username asc").
		Find(&records)

	users := make([]interfaces.StatsUserStruct, 0, len(records))
	for _, record := range records {
		users = append(users, interfaces.StatsUserStruct{
			ID:       record.UserID,
			Username: record.Username,
			Posts:    record.Posts,
			LastPost: record.LastPost,
		})
	}

	return users, nil
}

func (r UserStatsRepo) GetTopUsers() ([]interfaces.StatsUserStruct, error) {
	var records []interfaces.Stats

	r.tx.
		Where("user_id != 0").
		Order("posts desc").
		Find(&records)

	users := make([]interfaces.StatsUserStruct, 0, len(records))
	for _, record := range records {
		users = append(users, interfaces.StatsUserStruct{
			ID:       record.UserID,
			Username: record.Username,
			Posts:    record.Posts,
			LastPost: record.LastPost,
		})
	}

	return users, nil
}
