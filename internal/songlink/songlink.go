package songlink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type ResponseEntity struct {
	Title  string `json:"title"`
	Artist string `json:"artistName"`
}

type ResponseLink struct {
	URL string `json:"url"`
}

type ResponseLinks struct {
	Spotify    ResponseLink
	AppleMusic ResponseLink
	Youtube    ResponseLink
}

type Response struct {
	PageURL  string                    `json:"pageUrl"`
	Entities map[string]ResponseEntity `json:"entitiesByUniqueId"`
	Links    ResponseLinks             `json:"linksByPlatform"`
}

type Entry struct {
	Type   EntryType
	Title  string
	Artist string
	Links  []EntryLink
}

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

type EntryLink struct {
	Platform Platform
	URL      string
}

type Platform string

const (
	Songlink   Platform = "songlink"
	Spotify    Platform = "spotify"
	AppleMusic Platform = "appleMusic"
	Youtube    Platform = "youtube"
)

func (p Platform) Natural() string {
	switch p {
	case Songlink:
		return "Songlink"
	case Spotify:
		return "Spotify"
	case AppleMusic:
		return "Apple Music"
	case Youtube:
		return "YouTube"
	default:
		return "Unknown"
	}
}

func GetSonglinkEntry(url string) (*Entry, error) {
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
		platform = Spotify
		_type = Song
		ID = matches[1]
	} else if spotifyAlbumRegex.MatchString(url) {
		matches := spotifyAlbumRegex.FindStringSubmatch(url)
		platform = Spotify
		_type = Album
		ID = matches[1]
	} else if appleMusicTrackRegex.MatchString(url) {
		matches := appleMusicTrackRegex.FindStringSubmatch(url)
		platform = AppleMusic
		_type = Song
		ID = matches[1]
	} else if appleMusicAlbumRegex.MatchString(url) {
		matches := appleMusicAlbumRegex.FindStringSubmatch(url)
		platform = AppleMusic
		_type = Album
		ID = matches[1]
	} else {
		return nil, fmt.Errorf("unable to parse platform, type and ID from %s", url)
	}

	return getSonglinkEntry(platform, _type, ID)
}

func getSonglinkEntry(platform Platform, _type EntryType, ID string) (*Entry, error) {
	entry := Entry{
		Type:  _type,
		Links: make([]EntryLink, 0),
	}

	// For all other data we need to do an request to the Songlink API
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
		Platform: Songlink,
		URL:      response.PageURL,
	})
	if platform != Spotify && response.Links.Spotify.URL != "" {
		entry.Links = append(entry.Links, EntryLink{
			Platform: Spotify,
			URL:      response.Links.Spotify.URL,
		})
	}
	if platform != AppleMusic && response.Links.AppleMusic.URL != "" {
		entry.Links = append(entry.Links, EntryLink{
			Platform: AppleMusic,
			URL:      response.Links.AppleMusic.URL,
		})
	}
	if response.Links.Youtube.URL != "" {
		entry.Links = append(entry.Links, EntryLink{
			Platform: Youtube,
			URL:      response.Links.Youtube.URL,
		})
	}

	return &entry, nil
}

// Example URLs:
//
// Spotify album: https://api.song.link/v1-alpha.1/links?platform=spotify&type=album&id=2Gbv0Wjtwn9zQYMvWtTHnK
// Spotify track: https://api.song.link/v1-alpha.1/links?platform=spotify&type=song&id=0Q5IOvNoREy7gzT0CWmayo
// Apple Music album: https://api.song.link/v1-alpha.1/links?platform=appleMusic&type=album&id=1472283462
// Apple Music track: https://api.song.link/v1-alpha.1/links?platform=appleMusic&type=song&id=1472283463

func getSonglinkResponse(platform Platform, _type EntryType, ID string) (*Response, error) {
	url := fmt.Sprintf(
		"https://api.song.link/v1-alpha.1/links?platform=%s&type=%s&id=%s&userCountry=DE",
		platform,
		_type,
		ID,
	)

	request, _ := http.NewRequest("GET", url, nil)
	response, _ := http.DefaultClient.Do(request)

	responseBody := &Response{}
	if err := json.NewDecoder(response.Body).Decode(responseBody); err != nil {
		return nil, err
	}

	return responseBody, nil
}
