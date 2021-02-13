package config

import (
	"log"
	"strings"

	"github.com/neovg/kmptnzbot/internal/logger"
	"github.com/spf13/viper"
)

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
	BaseUrl             string
	EndpointSendMessage string
}

type StonksMatcher struct {
	QuotesUrl string `yaml:"quotesUrl"`
}

type BuzzwordsMatcher struct {
	Trigger string `yaml:"trigger"`
	Reply   string `yaml:"reply"`
}

var Cfg Config

func init() {
	logger.Log.Debug("init config")

	// Search config files in current directory only
	viper.AddConfigPath(".")

	// Load config file
	viper.SetConfigFile("config.yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicln(err)
	}

	// Load .env file
	viper.SetConfigFile(".env")
	if err := viper.MergeInConfig(); err != nil {
		log.Println("no .env file found")
	}

	// Mapping between keys in .env file or environment to config
	envToConfig := map[string]string{
		"listen_addr":       "server.listenAddr",
		"db_driver":         "database.driver",
		"sqlite_file":       "database.sqlite.file",
		"postgres_host":     "database.postgres.host",
		"postgres_port":     "database.postgres.port",
		"postgres_dbname":   "database.postgres.dbname",
		"postgres_user":     "database.postgres.user",
		"postgres_password": "database.postgres.password",
		"postgres_ssl":      "database.postgres.ssl",
		"postgres_timezone": "database.postgres.timezone",
		"telegram_api_key":  "telegram.apiKey",
	}

	// Map directives from environment variables to config
	for envKey, configKey := range envToConfig {
		// Value from .env file overwrites value from config.yml
		val := viper.GetString(envKey)
		if len(val) > 0 {
			viper.Set(configKey, val)
		}

		// Bind environment variable to config
		if err := viper.BindEnv(configKey, strings.ToUpper(envKey)); err != nil {
			log.Panicln(err)
		}
	}

	// Convert completed config data in Viper to Config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panicln(err)
	}

	Cfg = cfg
}
