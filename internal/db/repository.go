package db

import (
	"github.com/br0-space/bot/interfaces"
	"gorm.io/gorm"
)

type Repository struct {
	log          interfaces.LoggerInterface
	messageStats interfaces.MessageStatsRepoInterface
	plusplus     interfaces.PlusplusRepoInterface
	stats        interfaces.StatsRepoInterface
}

func NewRepository(
	logger interfaces.LoggerInterface,
	messageStats interfaces.MessageStatsRepoInterface,
	plusplus interfaces.PlusplusRepoInterface,
	stats interfaces.StatsRepoInterface,
) *Repository {
	return &Repository{
		log:          logger,
		messageStats: messageStats,
		plusplus:     plusplus,
		stats:        stats,
	}
}

func (r *Repository) MessageStats() interfaces.MessageStatsRepoInterface {
	return r.messageStats
}

func (r *Repository) Plusplus() interfaces.PlusplusRepoInterface {
	return r.plusplus
}

func (r *Repository) Stats() interfaces.StatsRepoInterface {
	return r.stats
}

type BaseRepo struct {
	log   interfaces.LoggerInterface
	model interface{}
	tx    *gorm.DB
}

func NewBaseRepo(logger interfaces.LoggerInterface, model interface{}, tx *gorm.DB) BaseRepo {
	return BaseRepo{
		log:   logger,
		model: model,
		tx:    tx,
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
