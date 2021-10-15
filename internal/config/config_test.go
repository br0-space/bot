package config

import (
	"testing"

	_ "github.com/br0-space/bot/testing_init"

	"github.com/stretchr/testify/assert"
)

func TestNewTestConfig(t *testing.T) {
	cfg := NewTestConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, cfg.Server.ListenAddr, ":3000")
	assert.Equal(t, cfg.Database.Driver, "sqlite")
}
