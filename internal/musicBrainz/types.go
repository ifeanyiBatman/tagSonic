package musicbrainz

// mbRecordingResponse maps the JSON from the MusicBrainz Recording API.
// Endpoint: /ws/2/recording/{mbid}?inc=artist-credits+releases+genres&fmt=json
type mbRecordingResponse struct {
	Title        string `json:"title"`
	ArtistCredit []struct {
		Name   string `json:"name"`
		Artist struct {
			Name string `json:"name"`
		} `json:"artist"`
	} `json:"artist-credit"`
	Releases []struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		Date         string `json:"date"`
		ReleaseGroup struct {
			PrimaryType string `json:"primary-type"`
		} `json:"release-group"`
	} `json:"releases"`
	Genres []struct {
		Name string `json:"name"`
	} `json:"genres"`
}

// SongTags holds the cleaned-up metadata extracted from MusicBrainz.
type SongTags struct {
	Title       string
	Artist      string
	Album       string
	Genre       string
	Year        string
	CoverArtURL string
}
