package config

import (
	"log"

	"github.com/spf13/viper"
)

var testConfig *Config

func NewTestConfig() *Config {
	if testConfig == nil {
		testConfig = readTestConfig()
	}
	return testConfig
}

func readTestConfig() *Config {
	// Search config files in current directory only
	viper.AddConfigPath(".")

	// Read from config template file only
	viper.SetConfigFile("config.yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicln(err)
	}

	// Convert completed config data in Viper to Config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panicln(err)
	}

	return &cfg
}
