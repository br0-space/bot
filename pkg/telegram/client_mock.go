package telegram

import "github.com/br0-space/bot/interfaces"

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (c MockClient) SendMessage(chatID int64, messageOut interfaces.TelegramMessageStruct) error {
	return nil
}
