package db

import (
	"github.com/br0-space/bot/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
)

var mutexPlusplus sync.Mutex

type PlusplusRepo struct {
	BaseRepo
	log interfaces.LoggerInterface
	tx  *gorm.DB
}

func NewPlusplusRepo(logger interfaces.LoggerInterface, tx *gorm.DB) *PlusplusRepo {
	return &PlusplusRepo{
		BaseRepo: NewBaseRepo(
			logger,
			&Plusplus{},
			tx,
		),
		log: logger,
		tx:  tx,
	}
}

func (r *PlusplusRepo) Increment(chatID int64, name string, increment int) (int, error) {
	mutexPlusplus.Lock()
	defer mutexPlusplus.Unlock()

	if err := r.tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "chat_id"}, {Name: "name"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"value": gorm.Expr("plusplus.value + ?", increment),
		}),
	}).Create(&Plusplus{
		ChatID: chatID,
		Name:   name,
		Value:  increment,
	}).Error; err != nil {
		return 0, err
	}

	var record Plusplus
	if err := r.tx.Where("chat_id = ? AND name = ?", chatID, name).First(&record).Error; err != nil {
		return 0, err
	}

	return record.Value, nil
}
