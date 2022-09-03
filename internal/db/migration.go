package db

import (
	"github.com/br0-space/bot/interfaces"
)

type DatabaseMigration struct {
	log  interfaces.LoggerInterface
	repo interfaces.DatabaseRepositoryInterface
}

func NewDatabaseMigration(
	log interfaces.LoggerInterface,
	repo interfaces.DatabaseRepositoryInterface,
) *DatabaseMigration {
	return &DatabaseMigration{
		log:  log,
		repo: repo,
	}
}

func (m *DatabaseMigration) Migrate() error {
	m.log.Debug("Migrating table", m.repo.MessageStats().(interfaces.RepoInterface).TableName())
	if err := m.repo.MessageStats().(interfaces.RepoInterface).Migrate(); err != nil {
		return err
	}

	m.log.Debug("Migrating table", m.repo.Plusplus().(interfaces.RepoInterface).TableName())
	if err := m.repo.Plusplus().(interfaces.RepoInterface).Migrate(); err != nil {
		return err
	}

	m.log.Debug("Migrating table", m.repo.Stats().(interfaces.RepoInterface).TableName())
	if err := m.repo.Stats().(interfaces.RepoInterface).Migrate(); err != nil {
		return err
	}

	return nil
}
