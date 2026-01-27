package repo

import (
	"sync"

	"github.com/br0-space/bot/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var mutexPlusplus sync.Mutex

type PlusplusRepo struct {
	BaseRepo
}

func NewPlusplusRepo(tx *gorm.DB) *PlusplusRepo {
	return &PlusplusRepo{
		BaseRepo: NewBaseRepo(
			tx,
			&interfaces.Plusplus{
				Name:  "",
				Value: 0,
			},
		),
	}
}

func (r PlusplusRepo) Increment(name string, increment int) (int, error) {
	mutexPlusplus.Lock()
	defer mutexPlusplus.Unlock()

	if err := r.tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.Assignments(map[string]any{
			"value": gorm.Expr("plusplus.value + ?", increment),
		}),
	}).Create(&interfaces.Plusplus{
		Name:  name,
		Value: increment,
	}).Error; err != nil {
		return 0, err
	}

	var record interfaces.Plusplus
	if err := r.tx.
		Where("name = ?", name).
		First(&record).
		Error; err != nil {
		return 0, err
	}

	return record.Value, nil
}

func (r PlusplusRepo) FindTops(limit int) ([]interfaces.Plusplus, error) {
	var records []interfaces.Plusplus
	if err := r.tx.
		Order("value desc").
		Limit(limit).
		Find(&records).
		Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r PlusplusRepo) FindFlops(limit int) ([]interfaces.Plusplus, error) {
	var records []interfaces.Plusplus
	if err := r.tx.
		Order("value asc").
		Limit(limit).
		Find(&records).
		Error; err != nil {
		return nil, err
	}

	return records, nil
}
