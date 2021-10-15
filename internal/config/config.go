package config

type Config struct {
	Server           Server
	Database         Database
	Telegram         Telegram
	StonksMatcher    StonksMatcher
	BuzzwordsMatcher []BuzzwordsMatcher
}

type Server struct {
	ListenAddr string
}

type Database struct {
	Driver string
	SQLite struct {
		File string
	}
	Postgres struct {
		Host     string
		Port     uint
		DBName   string
		User     string
		Password string
		SSL      string
		Timezone string
	}
}

type Telegram struct {
	ApiKey              string
	WebhookURL          string
	BaseUrl             string
	EndpointSetWebhook  string
	EndpointSendMessage string
}

type StonksMatcher struct {
	QuotesUrl string
}

type BuzzwordsMatcher struct {
	Trigger string
	Reply   string
}