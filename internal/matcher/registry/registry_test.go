package registry

import (
	"testing"

	"github.com/br0-space/bot/internal/telegram/webhook"
	"github.com/stretchr/testify/assert"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	assert.NotNil(t, registry)
}

func TestRegistry_ProcessWebhookMessageInAllMatchers(t *testing.T) {
	message := webhook.MakeTestMessage(1, 2, "foobar")
	registry := NewRegistry()
	registry.ProcessWebhookMessageInAllMatchers(message)
}