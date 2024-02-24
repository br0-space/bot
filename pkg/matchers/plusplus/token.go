package plusplus

import (
	"fmt"
	"strconv"

	telegramclient "github.com/br0-space/bot-telegramclient"
)

type Token struct {
	Name      string
	Increment int
}

func (t Token) MakeReply(value int) telegramclient.MessageStruct {
	mode := strconv.Itoa(t.Increment)

	if t.Increment > 0 {
		mode = "+" + mode
	}

	return telegramclient.MarkdownMessage(
		fmt.Sprintf(
			template,
			telegramclient.EscapeMarkdown(mode),
			telegramclient.EscapeMarkdown(t.Name),
			telegramclient.EscapeMarkdown(strconv.Itoa(value)),
		),
	)
}
