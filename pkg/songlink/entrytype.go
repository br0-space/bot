package songlink

type EntryType string

const (
	Song  EntryType = "song"
	Album EntryType = "album"
)

func (t EntryType) Natural() string {
	switch t {
	case Song:
		return "Song"
	case Album:
		return "Album"
	default:
		return "Unknown"
	}
}
