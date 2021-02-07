package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/neovg/kmptnzbot/internal/logger"
	"gopkg.in/yaml.v2"
)

type Server struct {
	ListenAddr string `yaml:"listenAddr"`
}

type Telegram struct {
	ApiKey              string `yaml:"apiKey"`
	BaseUrl             string `yaml:"baseUrl"`
	EndpointSendMessage string `yaml:"endpointSendMessage"`
}

type StonksMatcher struct {
	QuotesUrl string `yaml:"quotesUrl"`
}

type Config struct {
	Server        Server        `yaml:"server"`
	Telegram      Telegram      `yaml:"telegram"`
	StonksMatcher StonksMatcher `yaml:"stonksMatcher"`
}

var Cfg *Config

func init() {
	logger.Log.Debug("init config")

	// Generate our config based on the config supplied
	// by the user in the flags
	cfgPath, err := ParseFlags()
	if err != nil {
		logger.Log.Fatal(err)
	}
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		logger.Log.Fatal(err)
	}

	Cfg = cfg
}

// NewConfig returns a new decoded Config struct
func NewConfig(cfgPath string) (*Config, error) {
	// Create config structure
	cfg := &Config{}

	// Open config file
	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&cfg); err != nil {
		err = file.Close()
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}
