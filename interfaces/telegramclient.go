package interfaces

type TelegramClientInterface interface {
	SendMessage(chatID int64, messageOut TelegramMessageStruct) error
}

// Create a struct that is accepted by Telegram's sendMessage endpoint
// https://core.telegram.org/bots/api#sendmessage

type TelegramMessageStruct struct {
	ChatID                int64  `json:"chat_id"`
	ReplyToMessageID      int64  `json:"reply_to_message_id"`
	Text                  string `json:"text"`
	Photo                 string `json:"photo"`
	Caption               string `json:"caption"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotification   bool   `json:"disable_notification"`
}

type TelegramMessageResponseBodyStruct struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int16  `json:"error_code"`
	Description string `json:"description"`
}
