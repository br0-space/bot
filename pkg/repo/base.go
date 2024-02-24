package repo

import (
	logger "github.com/br0-space/bot-logger"
	"gorm.io/gorm"
)

type BaseRepo struct {
	log   logger.Interface
	tx    *gorm.DB
	model interface{}
}

func NewBaseRepo(tx *gorm.DB, model interface{}) BaseRepo {
	return BaseRepo{
		log:   logger.New(),
		tx:    tx,
		model: model,
	}
}

func (r *BaseRepo) Model() interface{} {
	return r.model
}

func (r *BaseRepo) TableName() string {
	stmt := &gorm.Statement{DB: r.tx} //nolint:exhaustruct
	_ = stmt.Parse(r.Model())

	return stmt.Table
}

func (r *BaseRepo) Migrate() error {
	return r.tx.AutoMigrate(r.Model())
}
