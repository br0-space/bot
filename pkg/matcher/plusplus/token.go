package plusplus

import (
	"fmt"
	telegramclient "github.com/br0-space/bot-telegramclient"
)

type Token struct {
	Name      string
	Increment int
}

func (t Token) MakeReply(value int) telegramclient.MessageStruct {
	var mode string
	switch {
	case t.Increment > 0:
		mode = fmt.Sprintf("+%d", t.Increment)
	case t.Increment < 0:
		mode = fmt.Sprintf("%d", t.Increment)
	default:
		mode = "+-"
	}

	return telegramclient.MarkdownMessage(
		fmt.Sprintf(
			template,
			telegramclient.EscapeMarkdown(mode),
			telegramclient.EscapeMarkdown(t.Name),
			telegramclient.EscapeMarkdown(fmt.Sprintf("%d", value)),
		),
	)
}
