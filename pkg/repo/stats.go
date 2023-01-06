package repo

import (
	"github.com/br0-space/bot/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"
)

var mutexStats sync.Mutex

type StatsRepo struct {
	BaseRepo
	log interfaces.LoggerInterface
	tx  *gorm.DB
}

func NewStatsRepo(logger interfaces.LoggerInterface, tx *gorm.DB) *StatsRepo {
	return &StatsRepo{
		BaseRepo: NewBaseRepo(
			logger,
			&interfaces.Stats{},
			tx,
		),
		log: logger,
		tx:  tx,
	}
}

func (r *StatsRepo) UpdateStats(userID int64, username string) error {
	mutexStats.Lock()
	defer mutexStats.Unlock()

	return r.tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
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

func (r *StatsRepo) GetKnownUsers() ([]interfaces.StatsUserStruct, error) {
	var records []interfaces.Stats
	r.tx.
		Where("user_id != 0").
		Order("username asc").
		Find(&records)

	var users []interfaces.StatsUserStruct
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

func (r *StatsRepo) GetTopUsers() ([]interfaces.StatsUserStruct, error) {
	var records []interfaces.Stats
	r.tx.
		Where("user_id != 0").
		Order("posts desc").
		Find(&records)

	var users []interfaces.StatsUserStruct
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