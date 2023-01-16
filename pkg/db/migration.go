package db

import (
	"github.com/br0-space/bot/interfaces"
)

type DatabaseMigration struct {
	log              interfaces.LoggerInterface
	messageStatsRepo interfaces.MessageStatsRepoInterface
	plusplusRepo     interfaces.PlusplusRepoInterface
	userStatsRepo    interfaces.UserStatsRepoInterface
}

func MakeDatabaseMigration(
	logger interfaces.LoggerInterface,
	messageStatsRepo interfaces.MessageStatsRepoInterface,
	plusplusRepo interfaces.PlusplusRepoInterface,
	userStatsRepo interfaces.UserStatsRepoInterface,
) DatabaseMigration {
	return DatabaseMigration{
		log:              logger,
		messageStatsRepo: messageStatsRepo,
		plusplusRepo:     plusplusRepo,
		userStatsRepo:    userStatsRepo,
	}
}

func (m DatabaseMigration) Migrate() error {
	m.log.Debug("Migrating table", m.messageStatsRepo.(interfaces.RepoInterface).TableName())
	if err := m.messageStatsRepo.(interfaces.RepoInterface).Migrate(); err != nil {
		return err
	}

	m.log.Debug("Migrating table", m.plusplusRepo.(interfaces.RepoInterface).TableName())
	if err := m.plusplusRepo.(interfaces.RepoInterface).Migrate(); err != nil {
		return err
	}

	m.log.Debug("Migrating table", m.userStatsRepo.(interfaces.RepoInterface).TableName())
	if err := m.userStatsRepo.(interfaces.RepoInterface).Migrate(); err != nil {
		return err
	}

	return nil
}
