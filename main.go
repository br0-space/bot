package main

import (
	"math/rand"
	"time"

	"github.com/neovg/kmptnzbot/cmd"
)

func main() {
	// Seed rand before doing anything else
	rand.Seed(time.Now().UnixNano())

	// Run server
	cmd.RunServer()
}
