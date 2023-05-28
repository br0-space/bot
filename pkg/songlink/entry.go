package songlink

import (
	"fmt"

	telegramclient "github.com/br0-space/bot-telegramclient"
)

type Entry struct {
	Type   EntryType
	Title  string
	Artist string
	Links  []EntryLink
}

type EntryLink struct {
	Platform Platform
	URL      string
}

func (e Entry) ToMarkdown() string {
	text := fmt.Sprintf(
		"*%s*\n*%s* · %s\n\n",
		telegramclient.EscapeMarkdown(e.Title),
		telegramclient.EscapeMarkdown(e.Artist),
		e.Type.Natural(),
	)

	for i := range e.Links {
		if e.Links[i].Platform == PlatformSonglink {
			continue
		}

		text += fmt.Sprintf(
			"🎧 [%s](%s)\n\n",
			e.Links[i].Platform.Natural(),
			e.Links[i].URL,
		)
	}

	text += fmt.Sprintf(
		"🔗 [%s](%s)",
		PlatformSonglink.Natural(),
		e.Links[0].URL,
	)

	return text
}
