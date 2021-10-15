package telegram

type LiveTelegram struct {
	*telegram
}

func NewLiveTelegram() *LiveTelegram {
	return &LiveTelegram{
		telegram: newTelegram(),
	}
}

func (t *LiveTelegram) SendMessage(chatID int64, messageOut Message) error {
	t.telegram.SendMessage(chatID, messageOut)
	return nil
}

func (t *LiveTelegram) SentMessages() []Message {
	return t.telegram.SentMessages()
}