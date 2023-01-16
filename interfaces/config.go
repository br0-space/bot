package interfaces

type ConfigStruct struct {
	Verbose  bool `mapstructure:"verbose"`
	Quiet    bool `mapstructure:"quiet"`
	Server   ServerConfigStruct
	Database DatabaseConfigStruct
	Telegram TelegramConfigStruct
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

type TelegramConfigStruct struct {
	ApiKey              string
	WebhookURL          string
	BaseUrl             string
	EndpointSetWebhook  string
	EndpointSendMessage string
	ChatID              int64
}

type MatcherConfigStruct struct {
	Enabled bool
}
