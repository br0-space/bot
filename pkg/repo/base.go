package repo

import (
	"github.com/br0-space/bot/interfaces"
	"gorm.io/gorm"
)

type BaseRepo struct {
	log   interfaces.LoggerInterface
	tx    *gorm.DB
	model interface{}
}

func NewBaseRepo(logger interfaces.LoggerInterface, tx *gorm.DB, model interface{}) BaseRepo {
	return BaseRepo{
		log:   logger,
		tx:    tx,
		model: model,
	}
}

func (r *BaseRepo) Model() interface{} {
	return r.model
}

func (r *BaseRepo) TableName() string {
	stmt := &gorm.Statement{DB: r.tx}
	_ = stmt.Parse(r.Model())
	return stmt.Table
}

func (r *BaseRepo) Migrate() error {
	return r.tx.AutoMigrate(r.Model())
}
