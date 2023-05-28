package interfaces

import telegramclient "github.com/br0-space/bot-telegramclient"

type ConfigStruct struct {
	Verbose  bool `mapstructure:"verbose"`
	Quiet    bool `mapstructure:"quiet"`
	Server   ServerConfigStruct
	Database DatabaseConfigStruct
	Telegram telegramclient.ConfigStruct
}

type ServerConfigStruct struct {
	ListenAddr string
}

type DatabaseConfigStruct struct {
	Driver string
	SQLite struct {
		File string
	}
	PostgreSQL struct {
		Host     string
		Port     uint
		DBName   string
		User     string
		Password string
		SSL      string
		Timezone string
	}
	AutoMigrate bool
}

type MatcherConfigStruct struct {
	Enabled bool
}
