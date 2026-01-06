package interfaces

type DatabaseMigrationInterface interface {
	Migrate() error
}

type DatabaseRepositoryInterface interface {
	MessageStats() MessageStatsRepoInterface
	Plusplus() PlusplusRepoInterface
	Stats() UserStatsRepoInterface
}

type RepoInterface interface {
	Model() any
	TableName() string
	Migrate() error
}
