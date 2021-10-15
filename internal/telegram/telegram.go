package telegram

type Message struct {
	ChatID                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotification   bool   `json:"disable_notification"`
	ReplyToMessageID      int64  `json:"reply_to_message_id"`
}

type Interface interface {
	SendMessage(chatID int64, message Message) error
	SentMessages() []Message
}

type telegram struct {
	sentMessages []Message
}

func newTelegram() *telegram {
	return &telegram{
		sentMessages: make([]Message, 0),
	}
}

func (t *telegram) SendMessage(chatID int64, messageOut Message) Message {
	messageOut.ChatID = chatID
	t.sentMessages = append(t.sentMessages, messageOut)
	return messageOut
}

func (t *telegram) SentMessages() []Message {
	return t.sentMessages
}