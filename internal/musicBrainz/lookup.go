package musicbrainz

import (
	"fmt"

	"github.com/michiwend/gomusicbrainz"
)

// LookupByMBID queries the MusicBrainz API using a Recording MBID (from AcoustID)
// and returns a SongTags struct with the full metadata.
// gomusicbrainz handles rate limiting and User-Agent automatically.
func LookupByMBID(recordingID string) (*SongTags, error) {
	client, err := gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"tagSonic",
		"1.0.0",
		"https://github.com/ifeanyiBatman/tagSonic",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MusicBrainz client: %w", err)
	}

	// LookupRecording fetches a recording by MBID with artist-credits included.
	// The gomusicbrainz Recording struct includes Title and ArtistCredit.
	recording, err := client.LookupRecording(gomusicbrainz.MBID(recordingID), "artist-credits", "releases")
	if err != nil {
		return nil, fmt.Errorf("MusicBrainz lookup failed for %s: %w", recordingID, err)
	}

	tags := &SongTags{
		Title: recording.Title,
	}

	// Extract artist from artist credits
	if len(recording.ArtistCredit.NameCredits) > 0 {
		tags.Artist = recording.ArtistCredit.NameCredits[0].Artist.Name
	}

	return tags, nil
}
