package music

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

type SonglinkEntry struct {
	Type          string
	Title         string
	Artist        string
	SonglinkURL   string
	SpotifyURL    string
	AppleMusicURL string
	YoutubeURL    string
}

// Examples URLs:
//
// Spotify track: https://open.spotify.com/track/0Q5IOvNoREy7gzT0CWmayo?si=d2b1a4b4ae204358
// Spotify album: https://open.spotify.com/album/2Gbv0Wjtwn9zQYMvWtTHnK
// Apple Music track 1: https://music.apple.com/de/album/by-1899-the-age-of-outlaws-and-gunslingers-was-at-an-end/1472283462?i=1472283463
// Apple Music track 2: https://music.apple.com/de/album/hi/1140071785?i=1140071869&l=en
// Apple Music album 1: https://music.apple.com/de/album/the-music-of-red-dead-redemption-2-original-score/1472283462
// Apple Music album 2: https://music.apple.com/de/album/life-thrills/1140071785?l=en
func GetSonglinkEntry(url string) (*SonglinkEntry, error) {
	const spotifyTrackPattern = "https:\\/\\/open\\.spotify\\.com\\/track\\/([a-zA-Z0-9]+)"
	const spotifyAlbumPattern = "https:\\/\\/open\\.spotify\\.com\\/album\\/([a-zA-Z0-9]+)"
	const appleMusicTrackPattern = "https:\\/\\/music\\.apple\\.com\\/[a-z]{2}\\/album\\/.+?\\/[0-9]+\\?i=([0-9]+)"
	const appleMusicAlbumPattern = "https:\\/\\/music\\.apple\\.com\\/[a-z]{2}\\/album\\/.+?\\/([0-9]+)"

	spotifyAlbumRegex := regexp.MustCompile(spotifyAlbumPattern)
	spotifyTrackRegex := regexp.MustCompile(spotifyTrackPattern)
	appleMusicAlbumRegex := regexp.MustCompile(appleMusicAlbumPattern)
	appleMusicTrackRegex := regexp.MustCompile(appleMusicTrackPattern)

	var platform string
	var _type string
	var ID string

	if spotifyTrackRegex.MatchString(url) {
		matches := spotifyTrackRegex.FindStringSubmatch(url)
		platform = "spotify"
		_type = "song"
		ID = matches[1]
	} else if spotifyAlbumRegex.MatchString(url) {
		matches := spotifyAlbumRegex.FindStringSubmatch(url)
		platform = "spotify"
		_type = "album"
		ID = matches[1]
	} else if appleMusicTrackRegex.MatchString(url) {
		matches := appleMusicTrackRegex.FindStringSubmatch(url)
		platform = "appleMusic"
		_type = "song"
		ID = matches[1]
	} else if appleMusicAlbumRegex.MatchString(url) {
		matches := appleMusicAlbumRegex.FindStringSubmatch(url)
		platform = "appleMusic"
		_type = "album"
		ID = matches[1]
	} else {
		return nil, fmt.Errorf("unable to parse platform, type and ID from %s", url)
	}

	return getSonglinkEntry(platform, _type, ID)
}

func getSonglinkEntry(platform string, _type string, ID string) (*SonglinkEntry, error) {
	entry := SonglinkEntry{}

	// The type is already known, so we can immediately add it to entry
	if _type == "album" {
		entry.Type = "Album"
	} else if _type == "song" {
		entry.Type = "Track"
	}

	// For all other data we need to do an request to the Songlink API
	response, err := getSonglinkResponse(platform, _type, ID)
	if err != nil {
		return nil, err
	}

	// Since the original URL came from a certain platform, an entity from that platform must exist,
	// so we can take title and artist from there
	var entityKey string
	if platform == "spotify" && _type == "album" {
		entityKey = fmt.Sprintf("SPOTIFY_ALBUM::%s", ID)
	} else if platform == "spotify" && _type == "song" {
		entityKey = fmt.Sprintf("SPOTIFY_SONG::%s", ID)
	} else if platform == "appleMusic" && _type == "album" {
		entityKey = fmt.Sprintf("ITUNES_ALBUM::%s", ID)
	} else if platform == "appleMusic" && _type == "song" {
		entityKey = fmt.Sprintf("ITUNES_SONG::%s", ID)
	}
	entry.Title = response.Entities[entityKey].Title
	entry.Artist = response.Entities[entityKey].Artist

	// Now we add all links
	entry.SonglinkURL = response.PageURL
	entry.SpotifyURL = response.Links.Spotify.URL
	entry.AppleMusicURL = response.Links.AppleMusic.URL
	entry.YoutubeURL = response.Links.Youtube.URL

	return &entry, nil
}

// Example URLs:
//
// Spotify album: https://api.song.link/v1-alpha.1/links?platform=spotify&type=album&id=2Gbv0Wjtwn9zQYMvWtTHnK
// Spotify track: https://api.song.link/v1-alpha.1/links?platform=spotify&type=song&id=0Q5IOvNoREy7gzT0CWmayo
// Apple Music album: https://api.song.link/v1-alpha.1/links?platform=appleMusic&type=album&id=1472283462
// Apple Music track: https://api.song.link/v1-alpha.1/links?platform=appleMusic&type=song&id=1472283463
func getSonglinkResponse(platform string, _type string, ID string) (*Response, error) {
	url := fmt.Sprintf("https://api.song.link/v1-alpha.1/links?platform=%s&type=%s&id=%s&userCountry=DE", platform, _type, ID)

	request, _ := http.NewRequest("GET", url, nil)
	response, _ := http.DefaultClient.Do(request)

	responseBody := &Response{}
	if err := json.NewDecoder(response.Body).Decode(responseBody); err != nil {
		return nil, err
	}

	return responseBody, nil
}
