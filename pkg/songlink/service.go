package songlink

import (
	"fmt"
	"github.com/br0-space/bot/interfaces"
	"regexp"
)

type Service struct{}

func MakeService() interfaces.SonglinkServiceInterface {
	return Service{}
}

func (s Service) GetEntryForUrl(url string) (interfaces.SonglinkEntryInterface, error) {
	const spotifyTrackPattern = "https:\\/\\/open\\.spotify\\.com\\/track\\/([a-zA-Z0-9]+)"
	const spotifyAlbumPattern = "https:\\/\\/open\\.spotify\\.com\\/album\\/([a-zA-Z0-9]+)"
	const appleMusicTrackPattern = "https:\\/\\/music\\.apple\\.com\\/[a-z]{2}\\/album\\/.+?\\/[0-9]+\\?i=([0-9]+)"
	const appleMusicAlbumPattern = "https:\\/\\/music\\.apple\\.com\\/[a-z]{2}\\/album\\/.+?\\/([0-9]+)"

	spotifyAlbumRegex := regexp.MustCompile(spotifyAlbumPattern)
	spotifyTrackRegex := regexp.MustCompile(spotifyTrackPattern)
	appleMusicAlbumRegex := regexp.MustCompile(appleMusicAlbumPattern)
	appleMusicTrackRegex := regexp.MustCompile(appleMusicTrackPattern)

	var platform Platform
	var _type EntryType
	var ID string

	if spotifyTrackRegex.MatchString(url) {
		matches := spotifyTrackRegex.FindStringSubmatch(url)
		platform = PlatformSpotify
		_type = Song
		ID = matches[1]
	} else if spotifyAlbumRegex.MatchString(url) {
		matches := spotifyAlbumRegex.FindStringSubmatch(url)
		platform = PlatformSpotify
		_type = Album
		ID = matches[1]
	} else if appleMusicTrackRegex.MatchString(url) {
		matches := appleMusicTrackRegex.FindStringSubmatch(url)
		platform = PlatformAppleMusic
		_type = Song
		ID = matches[1]
	} else if appleMusicAlbumRegex.MatchString(url) {
		matches := appleMusicAlbumRegex.FindStringSubmatch(url)
		platform = PlatformAppleMusic
		_type = Album
		ID = matches[1]
	} else {
		return nil, fmt.Errorf("unable to parse platform, type and ID from %s", url)
	}

	return newSonglinkEntry(platform, _type, ID)
}

func newSonglinkEntry(platform Platform, _type EntryType, ID string) (interfaces.SonglinkEntryInterface, error) {
	entry := Entry{
		Type:  _type,
		Links: make([]EntryLink, 0),
	}

	// For all other data we need to do a request to the Songlink API
	response, err := getSonglinkResponse(platform, _type, ID)
	if err != nil {
		return nil, err
	}

	// Since the original URL came from a certain platform, an entity from that platform must exist,
	// so we can take title and artist from there
	var entityKey string
	if platform == "spotify" && _type == Album {
		entityKey = fmt.Sprintf("SPOTIFY_ALBUM::%s", ID)
	} else if platform == "spotify" && _type == Song {
		entityKey = fmt.Sprintf("SPOTIFY_SONG::%s", ID)
	} else if platform == "appleMusic" && _type == Album {
		entityKey = fmt.Sprintf("ITUNES_ALBUM::%s", ID)
	} else if platform == "appleMusic" && _type == Song {
		entityKey = fmt.Sprintf("ITUNES_SONG::%s", ID)
	}
	entry.Title = response.Entities[entityKey].Title
	entry.Artist = response.Entities[entityKey].Artist

	// Now we add all links
	entry.Links = append(entry.Links, EntryLink{
		Platform: PlatformSonglink,
		URL:      response.PageURL,
	})
	if platform != PlatformSpotify && response.Links.Spotify.URL != "" {
		entry.Links = append(entry.Links, EntryLink{
			Platform: PlatformSpotify,
			URL:      response.Links.Spotify.URL,
		})
	}
	if platform != PlatformAppleMusic && response.Links.AppleMusic.URL != "" {
		entry.Links = append(entry.Links, EntryLink{
			Platform: PlatformAppleMusic,
			URL:      response.Links.AppleMusic.URL,
		})
	}
	if response.Links.Youtube.URL != "" {
		entry.Links = append(entry.Links, EntryLink{
			Platform: PlatformYoutube,
			URL:      response.Links.Youtube.URL,
		})
	}

	return &entry, nil
}