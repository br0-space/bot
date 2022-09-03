package interfaces

type ConfigStruct struct {
	Verbose          bool `mapstructure:"verbose"`
	Quiet            bool `mapstructure:"quiet"`
	Server           ServerConfigStruct
	Database         DatabaseConfigStruct
	Telegram         TelegramConfigStruct
	Matchers         MatchersConfigStruct
	StonksMatcher    StonksMatcherConfigStruct
	BuzzwordsMatcher []BuzzwordsMatcherConfigStruct
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
}

type TelegramConfigStruct struct {
	ApiKey              string
	WebhookURL          string
	BaseUrl             string
	EndpointSetWebhook  string
	EndpointSendMessage string
	ChatID              int64
}

type MatchersConfigStruct struct {
	Atall      AtallMatcherConfigStruct
	Choose     ChooseMatcherConfigStruct
	Help       HelpMatcherConfigStruct
	Janein     JaneinMatcherConfigStruct
	Musiclinks MusicLinksMatcherConfigStruct
	Ping       PingMatcherConfigStruct
	Plusplus   PlusplusMatcherConfigStruct
	Stats      StatsMatcherConfigStruct
}

type AtallMatcherConfigStruct struct {
	Enabled bool
}

type ChooseMatcherConfigStruct struct {
	Enabled bool
}

type HelpMatcherConfigStruct struct {
	Enabled bool
}

type JaneinMatcherConfigStruct struct {
	Enabled bool
}

type MusicLinksMatcherConfigStruct struct {
	Enabled bool
}

type PingMatcherConfigStruct struct {
	Enabled bool
}

type PlusplusMatcherConfigStruct struct {
	Enabled bool
}

type StatsMatcherConfigStruct struct {
	Enabled bool
}

type StonksMatcherConfigStruct struct {
	QuotesUrl string
}

type BuzzwordsMatcherConfigStruct struct {
	Trigger string
	Reply   string
}
