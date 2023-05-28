package main

import (
	"github.com/br0-space/bot/container"
	"github.com/spf13/pflag"
)

func main() {
	pflag.Parse()

	logger := container.ProvideLogger()

	logger.Debug("debug")
	logger.Info("info")
	logger.Warning("warning")
	logger.Error("error")
	// logger.Panic("panic")
	// logger.Fatal("fatal")
}
