package telegram

type MockTelegram struct {
	*telegram
}

func NewMockTelegram() *MockTelegram {
	return &MockTelegram{
		telegram: newTelegram(),
	}
}

func (t *MockTelegram) SendMessage(chatID int64, messageOut Message) error {
	t.telegram.SendMessage(chatID, messageOut)
	return nil
}

func (t *MockTelegram) SentMessages() []Message {
	return t.telegram.SentMessages()
}