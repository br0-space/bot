package songlink

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiResponseEntity struct {
	Title  string `json:"title"`
	Artist string `json:"artistName"`
}

type apiResponseLink struct {
	URL string `json:"url"`
}

type apiResponseLinks struct { //nolint:musttag
	Spotify    apiResponseLink
	AppleMusic apiResponseLink
	Youtube    apiResponseLink
}

type apiResponse struct {
	PageURL  string                       `json:"pageUrl"`
	Entities map[string]apiResponseEntity `json:"entitiesByUniqueId"`
	Links    apiResponseLinks             `json:"linksByPlatform"`
}

// Example URLs:
//
// Spotify album: https://api.song.link/v1-alpha.1/links?platform=spotify&type=album&id=2Gbv0Wjtwn9zQYMvWtTHnK
// Spotify track: https://api.song.link/v1-alpha.1/links?platform=spotify&type=song&id=0Q5IOvNoREy7gzT0CWmayo
// Apple Music album: https://api.song.link/v1-alpha.1/links?platform=appleMusic&type=album&id=1472283462
// Apple Music track: https://api.song.link/v1-alpha.1/links?platform=appleMusic&type=song&id=1472283463

func getSonglinkResponse(platform Platform, _type EntryType, id string) (*apiResponse, error) {
	url := fmt.Sprintf(
		"https://api.song.link/v1-alpha.1/links?platform=%s&type=%s&id=%s&userCountry=DE",
		platform,
		_type,
		id,
	)

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	response, _ := http.DefaultClient.Do(request)

	responseBody := &apiResponse{}
	if err := json.NewDecoder(response.Body).Decode(responseBody); err != nil {
		return nil, err
	}

	return responseBody, nil
}
