package main

import (
	"math/rand"
	"time"

	"gitlab.com/br0-space/bot/cmd"
)

func main() {
	// Seed rand before doing anything else
	rand.Seed(time.Now().UnixNano())

	// Run server
	cmd.RunServer()
}
