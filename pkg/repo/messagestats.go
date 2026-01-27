package repo

import (
	"time"

	"github.com/br0-space/bot/interfaces"
	"gorm.io/gorm"
)

type MessageStatsRepo struct {
	BaseRepo
}

func NewMessageStatsRepo(tx *gorm.DB) *MessageStatsRepo {
	return &MessageStatsRepo{
		BaseRepo: NewBaseRepo(
			tx,
			&interfaces.MessageStats{
				UserID: 0,
				Time:   time.Time{},
				Words:  0,
			},
		),
	}
}

func (r MessageStatsRepo) InsertMessageStats(userID int64, words int) error {
	return r.tx.Create(&interfaces.MessageStats{
		UserID: userID,
		Time:   time.Now(),
		Words:  words,
	}).Error
}

func (r MessageStatsRepo) GetKnownUserIDs() ([]int64, error) {
	var userIDs []int64

	err := r.tx.
		Select("DISTINCT user_id").
		Where("user_id != 0").
		Find(&userIDs).
		Error

	return userIDs, err
}

func (r MessageStatsRepo) GetWordCounts() ([]interfaces.MessageStatsWordCountStruct, error) {
	var records []interfaces.MessageStatsWordCountStruct

	err := r.tx.Model(&interfaces.MessageStats{}).
		Joins("UserStats").
		Select(`"message_stats".user_id, "UserStats".username, count("message_stats".words) as words`).
		Where(`"message_stats"user_id != 0`).
		Group(`"message_stats".user_id, "UserStats".id, "UserStats".username`).
		Order(`count("message_stats".words) desc`).
		Scan(&records).
		Error

	return records, err
}
