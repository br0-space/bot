package config

import (
	"github.com/br0-space/bot/interfaces"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"strings"
	"time"
)

func init() {
	// Seed rand before doing anything else
	rand.Seed(time.Now().UnixNano())

	// Add default command line flags
	pflag.BoolP("verbose", "v", false, "Show verbose output")
	pflag.BoolP("quiet", "q", false, "Show errors only (overwrites verbose mode)")
}

func NewConfig() *interfaces.ConfigStruct {
	config, err := loadConfig()
	if err != nil {
		log.Panicln("Unable to load config:", err)
	}

	return config
}

func NewTestConfig() *interfaces.ConfigStruct {
	return &interfaces.ConfigStruct{
		Verbose:          false,
		Quiet:            false,
		Server:           interfaces.ServerConfigStruct{},
		Database:         interfaces.DatabaseConfigStruct{},
		Telegram:         interfaces.TelegramConfigStruct{},
		Matchers:         interfaces.MatchersConfigStruct{},
		StonksMatcher:    interfaces.StonksMatcherConfigStruct{},
		BuzzwordsMatcher: []interfaces.BuzzwordsMatcherConfigStruct{},
	}
}

func loadConfig() (*interfaces.ConfigStruct, error) {
	// Bind to command line flags
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, err
	}

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
		"postgres_host":        "database.postgresql.host",
		"postgres_port":        "database.postgresql.port",
		"postgres_db":          "database.postgresql.dbname",
		"postgres_user":        "database.postgresql.user",
		"postgres_password":    "database.postgresql.password",
		"postgres_ssl":         "database.postgresql.ssl",
		"postgres_timezone":    "database.postgresql.timezone",
		"db_automigrate":       "database.autoMigrate",
		"telegram_api_key":     "telegram.apiKey",
		"telegram_webhook_url": "telegram.webhookUrl",
		"telegram_chat_id":     "telegram.chatID",
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
	var cfg interfaces.ConfigStruct
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panicln(err)
	}

	return &cfg, nil
}
