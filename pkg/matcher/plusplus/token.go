package plusplus

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"github.com/br0-space/bot/pkg/telegram"
)

type Token struct {
	Name      string
	Increment int
}

func (t Token) MakeReply(value int) interfaces.TelegramMessageStruct {
	var mode string
	switch {
	case t.Increment > 0:
		mode = fmt.Sprintf("+%d", t.Increment)
	case t.Increment < 0:
		mode = fmt.Sprintf("%d", t.Increment)
	default:
		mode = "+-"
	}

	return telegram.MakeMarkdownMessage(
		fmt.Sprintf(
			template,
			telegram.EscapeMarkdown(mode),
			telegram.EscapeMarkdown(t.Name),
			telegram.EscapeMarkdown(fmt.Sprintf("%d", value)),
		),
	)
}
