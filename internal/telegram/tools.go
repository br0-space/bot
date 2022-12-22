package telegram

import (
	"github.com/br0-space/bot/interfaces"
	"regexp"
)

func NewMessage(text string) interfaces.TelegramMessageStruct {
	return interfaces.TelegramMessageStruct{
		Text: text,
	}
}

func NewMarkdownMessage(text string) interfaces.TelegramMessageStruct {
	return interfaces.TelegramMessageStruct{
		Text:      text,
		ParseMode: "MarkdownV2",
	}
}

func NewReply(text string, messageID int64) interfaces.TelegramMessageStruct {
	return interfaces.TelegramMessageStruct{
		Text:             text,
		ReplyToMessageID: messageID,
	}
}

func NewMarkdownReply(text string, messageID int64) interfaces.TelegramMessageStruct {
	return interfaces.TelegramMessageStruct{
		Text:             text,
		ReplyToMessageID: messageID,
		ParseMode:        "MarkdownV2",
	}
}

func EscapeMarkdown(text string) string {
	re := regexp.MustCompile("[_*\\[\\]()~`>#+\\-=|{}.!\\\\]")

	return re.ReplaceAllString(text, "\\$0")
}
