package xkcd

import (
	"fmt"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/nishanths/go-xkcd/v2"
)

const template = "*%s*\n\n_%s_\n\nxkcd [\\#%d](%s) \\(%d\\.%d\\.%d\\)"

type Comic struct {
	base xkcd.Comic
}

func FromComic(comic xkcd.Comic) Comic {
	return Comic{
		base: comic,
	}
}

func (c Comic) Number() int {
	return c.base.Number
}

func (c Comic) URL() string {
	return fmt.Sprintf(
		"https://xkcd.com/%d",
		c.base.Number,
	)
}

func (c Comic) ImageURL() string {
	return c.base.ImageURL
}

func (c Comic) ToMarkdown() string {
	return fmt.Sprintf(
		template,
		telegramclient.EscapeMarkdown(c.base.Title),
		telegramclient.EscapeMarkdown(c.base.Alt),
		c.base.Number,
		telegramclient.EscapeMarkdown(c.URL()),
		c.base.Day,
		c.base.Month,
		c.base.Year,
	)
}
