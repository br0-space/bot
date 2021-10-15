package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var liveConfig *Config

func NewLiveConfig() *Config {
	if liveConfig == nil {
		liveConfig = readLiveConfig()
	}
	return liveConfig
}

func readLiveConfig() *Config {
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
		"listen_addr":          "server.listenAddr",
		"db_driver":            "database.driver",
		"sqlite_file":          "database.sqlite.file",
		"postgres_host":        "database.postgres.host",
		"postgres_port":        "database.postgres.port",
		"postgres_db":          "database.postgres.dbname",
		"postgres_user":        "database.postgres.user",
		"postgres_password":    "database.postgres.password",
		"postgres_ssl":         "database.postgres.ssl",
		"postgres_timezone":    "database.postgres.timezone",
		"telegram_api_key":     "telegram.apiKey",
		"telegram_webhook_url": "telegram.webhookUrl",
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

	return &cfg
}