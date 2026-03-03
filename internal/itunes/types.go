package itunes

// itunesResponse maps the JSON response from the iTunes Search API.
type itunesResponse struct {
	ResultCount int            `json:"resultCount"`
	Results     []itunesResult `json:"results"`
}

type itunesResult struct {
	TrackName        string `json:"trackName"`
	ArtistName       string `json:"artistName"`
	CollectionName   string `json:"collectionName"`
	PrimaryGenreName string `json:"primaryGenreName"`
	ReleaseDate      string `json:"releaseDate"`
	ArtworkUrl100    string `json:"artworkUrl100"`
}

// SongTags holds the cleaned-up metadata extracted from the iTunes search result.
type SongTags struct {
	Title       string
	Artist      string
	Album       string
	Genre       string
	Year        string
	CoverArtURL string
}
