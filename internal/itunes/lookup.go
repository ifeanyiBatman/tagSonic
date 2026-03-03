package itunes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// LookupSong searches the iTunes Search API using a song title and artist name.
// It returns a SongTags struct with the metadata from the best match.
func LookupSong(title, artist string) (*SongTags, error) {
	searchTerm := url.QueryEscape(title + " " + artist)
	searchURL := fmt.Sprintf("https://itunes.apple.com/search?term=%s&entity=song&limit=5", searchTerm)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("iTunes search request failed: %w", err)
	}
	defer resp.Body.Close()

	var itunesResp itunesResponse
	if err := json.NewDecoder(resp.Body).Decode(&itunesResp); err != nil {
		return nil, fmt.Errorf("failed to decode iTunes response: %w", err)
	}

	if itunesResp.ResultCount == 0 || len(itunesResp.Results) == 0 {
		return nil, fmt.Errorf("no iTunes results found for: %s - %s", title, artist)
	}

	result := itunesResp.Results[0]

	// Extract the year from the release date (e.g. "2016-02-14T12:00:00Z" -> "2016")
	year := ""
	if len(result.ReleaseDate) >= 4 {
		year = result.ReleaseDate[:4]
	}

	// Upgrade the cover art URL from 100x100 to 1000x1000 for high resolution
	coverArtURL := strings.Replace(result.ArtworkUrl100, "100x100bb", "1000x1000bb", 1)

	return &SongTags{
		Title:       result.TrackName,
		Artist:      result.ArtistName,
		Album:       result.CollectionName,
		Genre:       result.PrimaryGenreName,
		Year:        year,
		CoverArtURL: coverArtURL,
	}, nil
}
