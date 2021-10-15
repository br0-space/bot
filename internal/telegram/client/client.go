package client

import (
	"github.com/br0-space/bot/container"
	"github.com/br0-space/bot/internal/config"
	"github.com/br0-space/bot/internal/logger"
	"github.com/br0-space/bot/internal/telegram"
)

type Client struct {
	log logger.Interface
	cfg config.Telegram
}

func NewClient() *Client {
	return &Client{
		log: container.ProvideLoggerService(),
		cfg: container.ProvideConfig().Telegram,
	}
}

func (c Client) sendMessage(messageOut telegram.Message) error {
	c.log.Debugf("send to %d: %s", messageOut.ChatID, messageOut.Text)
	return nil
}