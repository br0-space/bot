package main

import (
	"math/rand"
	"time"

	"github.com/br0-space/bot/cmd"
	"github.com/br0-space/bot/container"
)

func main() {
	logger := container.ProvideLoggerService()

	// Seed rand before doing anything else
	rand.Seed(time.Now().UnixNano())

	// Run server
	err := cmd.NewCmd().RunServer()
	if err != nil {
		logger.Fatal(err)
	}
}
