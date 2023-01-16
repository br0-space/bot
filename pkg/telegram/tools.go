package telegram

import (
	"github.com/br0-space/bot/interfaces"
	"regexp"
)

func MakeMessage(text string) interfaces.TelegramMessageStruct {
	return interfaces.TelegramMessageStruct{
		Text: text,
	}
}

func MakeMarkdownMessage(text string) interfaces.TelegramMessageStruct {
	return interfaces.TelegramMessageStruct{
		Text:      text,
		ParseMode: "MarkdownV2",
	}
}

func MakeReply(text string, messageID int64) interfaces.TelegramMessageStruct {
	return interfaces.TelegramMessageStruct{
		Text:             text,
		ReplyToMessageID: messageID,
	}
}

func MakeMarkdownReply(text string, messageID int64) interfaces.TelegramMessageStruct {
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
