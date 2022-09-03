package interfaces

import "time"

type DatabaseMigrationInterface interface {
	Migrate() error
}

type DatabaseRepositoryInterface interface {
	MessageStats() MessageStatsRepoInterface
	Plusplus() PlusplusRepoInterface
	Stats() StatsRepoInterface
}

type RepoInterface interface {
	Model() interface{}
	TableName() string
	Migrate() error
}

type MessageStatsRepoInterface interface {
	InsertMessageStats(chatID int64, userID int64, words int) error
	GetKnownUserIDs(chatID int64) ([]int64, error)
	GetWordCounts(chatID int64) ([]MessageStatsWordCountStruct, error)
}

type MessageStatsWordCountStruct struct {
	UserID   int64
	Username string
	Words    int
}

type PlusplusRepoInterface interface {
	Increment(chatID int64, name string, increment int) (int, error)
}

type StatsRepoInterface interface {
	UpdateStats(chatID int64, userID int64, username string) error
	GetKnownUsers(chatID int64) ([]StatsUserStruct, error)
	GetTopUsers(chatID int64) ([]StatsUserStruct, error)
}

type StatsUserStruct struct {
	ID       int64
	Username string
	Posts    uint32
	LastPost time.Time
}
