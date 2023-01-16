package main

import (
	"github.com/br0-space/bot/container"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/pflag"
)

func main() {
	pflag.Parse()

	logger := container.ProvideLogger()
	config := container.ProvideConfig()

	logger.Info("Config loaded from .env file, environment, and command line flags:")

	spew.Dump(config.Database)
}
