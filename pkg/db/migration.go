package db

import (
	logger "github.com/br0-space/bot-logger"
	"github.com/br0-space/bot/interfaces"
)

type DatabaseMigration struct {
	log              interfaces.LoggerInterface
	messageStatsRepo interfaces.MessageStatsRepoInterface
	plusplusRepo     interfaces.PlusplusRepoInterface
	rollRepo         interfaces.RollRepoInterface
	userStatsRepo    interfaces.UserStatsRepoInterface
}

func MakeDatabaseMigration(
	messageStatsRepo interfaces.MessageStatsRepoInterface,
	plusplusRepo interfaces.PlusplusRepoInterface,
	rollRepo interfaces.RollRepoInterface,
	userStatsRepo interfaces.UserStatsRepoInterface,
) DatabaseMigration {
	return DatabaseMigration{
		log:              logger.New(),
		messageStatsRepo: messageStatsRepo,
		plusplusRepo:     plusplusRepo,
		rollRepo:         rollRepo,
		userStatsRepo:    userStatsRepo,
	}
}

func (m DatabaseMigration) Migrate() error {
	if repo, ok := m.messageStatsRepo.(interfaces.RepoInterface); ok {
		m.log.Debug("Migrating table", repo.TableName())

		if err := repo.Migrate(); err != nil {
			return err
		}
	}

	if repo, ok := m.plusplusRepo.(interfaces.RepoInterface); ok {
		m.log.Debug("Migrating table", repo.TableName())

		if err := repo.Migrate(); err != nil {
			return err
		}
	}

	if repo, ok := m.rollRepo.(interfaces.RepoInterface); ok {
		m.log.Debug("Migrating table", repo.TableName())

		if err := repo.Migrate(); err != nil {
			return err
		}
	}

	if repo, ok := m.userStatsRepo.(interfaces.RepoInterface); ok {
		m.log.Debug("Migrating table", repo.TableName())

		if err := repo.Migrate(); err != nil {
			return err
		}
	}

	return nil
}
