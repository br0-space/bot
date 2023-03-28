package songlink_test

import (
	"github.com/br0-space/bot/pkg/songlink"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tests = []struct {
	in  string
	out *songlink.Entry
}{
	{
		in: "https://open.spotify.com/track/0Q5IOvNoREy7gzT0CWmayo?si=d2b1a4b4ae204358",
		out: &songlink.Entry{
			Type:   songlink.Song,
			Title:  "By 1899, The Age Of Outlaws And Gunslingers Was At An End",
			Artist: "Jeff Silverman, Luke O’Malley, Woody Jackson",
			Links: []songlink.EntryLink{
				{songlink.PlatformSonglink, "https://song.link/s/0Q5IOvNoREy7gzT0CWmayo"},
				{songlink.PlatformSpotify, "https://open.spotify.com/track/0Q5IOvNoREy7gzT0CWmayo"},
				{songlink.PlatformAppleMusic, "https://geo.music.apple.com/de/album/_/1472283462?i=1472283463&mt=1&app=music&ls=1&at=1000lHKX&ct=api_http&itscg=30200&itsct=odsl_m"},
				{songlink.PlatformYoutube, "https://www.youtube.com/watch?v=atmy8uAI8K0"},
			},
		},
	},
	{
		in: "https://open.spotify.com/album/2Gbv0Wjtwn9zQYMvWtTHnK",
		out: &songlink.Entry{
			Type:   songlink.Album,
			Title:  "The Music of Red Dead Redemption 2 (Original Score)",
			Artist: "Various Artists",
			Links: []songlink.EntryLink{
				{songlink.PlatformSonglink, "https://album.link/s/2Gbv0Wjtwn9zQYMvWtTHnK"},
				{songlink.PlatformSpotify, "https://open.spotify.com/album/2Gbv0Wjtwn9zQYMvWtTHnK"},
				{songlink.PlatformAppleMusic, "https://geo.music.apple.com/de/album/_/1472283462?mt=1&app=music&ls=1&at=1000lHKX&ct=api_http&itscg=30200&itsct=odsl_m"},
				{songlink.PlatformYoutube, "https://www.youtube.com/playlist?list=OLAK5uy_myT6DLJmO1jsviiIR4li7oyaHXWpyIVWo"},
			},
		},
	},
	{
		in: "https://music.apple.com/de/album/by-1899-the-age-of-outlaws-and-gunslingers-was-at-an-end/1472283462?i=1472283463",
		out: &songlink.Entry{
			Type:   songlink.Song,
			Title:  "By 1899, The Age of Outlaws and Gunslingers Was At an End",
			Artist: "Jeff Silverman, Luke O'Malley & Woody Jackson",
			Links: []songlink.EntryLink{
				{songlink.PlatformSonglink, "https://song.link/us/i/1472283463"},
				{songlink.PlatformSpotify, "https://open.spotify.com/track/0Q5IOvNoREy7gzT0CWmayo"},
				{songlink.PlatformAppleMusic, "https://geo.music.apple.com/de/album/_/1472283462?i=1472283463&mt=1&app=music&ls=1&at=1000lHKX&ct=api_http&itscg=30200&itsct=odsl_m"},
				{songlink.PlatformYoutube, "https://www.youtube.com/watch?v=atmy8uAI8K0"},
			},
		},
	},
	{
		in: "https://music.apple.com/de/album/hi/1140071785?i=1140071869&l=en",
		out: &songlink.Entry{
			Type:   songlink.Song,
			Title:  "Hi!",
			Artist: "Metrik",
			Links: []songlink.EntryLink{
				{songlink.PlatformSonglink, "https://song.link/us/i/1140071869"},
				{songlink.PlatformSpotify, "https://open.spotify.com/track/79lVzvbLCj9AXFGM7FRUVM"},
				{songlink.PlatformAppleMusic, "https://geo.music.apple.com/de/album/_/1140071785?i=1140071869&mt=1&app=music&ls=1&at=1000lHKX&ct=api_http&itscg=30200&itsct=odsl_m"},
				{songlink.PlatformYoutube, "https://www.youtube.com/watch?v=Q3enDjbwXWc"},
			},
		},
	},
	{
		in: "https://music.apple.com/de/album/the-music-of-red-dead-redemption-2-original-score/1472283462",
		out: &songlink.Entry{
			Type:   songlink.Album,
			Title:  "The Music of Red Dead Redemption 2 (Original Score)",
			Artist: "Verschiedene Interpreten",
			Links: []songlink.EntryLink{
				{songlink.PlatformSonglink, "https://album.link/us/i/1472283462"},
				{songlink.PlatformSpotify, "https://open.spotify.com/album/2Gbv0Wjtwn9zQYMvWtTHnK"},
				{songlink.PlatformAppleMusic, "https://geo.music.apple.com/de/album/_/1472283462?mt=1&app=music&ls=1&at=1000lHKX&ct=api_http&itscg=30200&itsct=odsl_m"},
			},
		},
	},
	{
		in: "https://music.apple.com/de/album/life-thrills/1140071785?l=en",
		out: &songlink.Entry{
			Type:   songlink.Album,
			Title:  "LIFE/THRILLS",
			Artist: "Metrik",
			Links: []songlink.EntryLink{
				{songlink.PlatformSonglink, "https://album.link/us/i/1140071785"},
				{songlink.PlatformSpotify, "https://open.spotify.com/album/3RTTmyLttk4sOdylJRDKrE"},
				{songlink.PlatformAppleMusic, "https://geo.music.apple.com/de/album/_/1140071785?mt=1&app=music&ls=1&at=1000lHKX&ct=api_http&itscg=30200&itsct=odsl_m"},
				{songlink.PlatformYoutube, "https://www.youtube.com/playlist?list=OLAK5uy_n3vkQMwLHzd3vClPzPEU9Oiy7COOwA89I"},
			},
		},
	},
}

func TestGetSonglinkEntry(t *testing.T) {
	t.Parallel()

	for _, tt := range tests {
		entry, err := songlink.MakeService().GetEntryForUrl(tt.in)
		assert.Nil(t, err)
		assert.NotNil(t, entry)
		assert.Equal(t, tt.out, entry)
	}
}
