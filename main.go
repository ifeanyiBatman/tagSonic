package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ifeanyiBatman/tagSonic/internal/acoustid"
	"github.com/ifeanyiBatman/tagSonic/internal/id3tagger"
	"github.com/ifeanyiBatman/tagSonic/internal/itunes"
	musicbrainz "github.com/ifeanyiBatman/tagSonic/internal/musicBrainz"
	"github.com/ifeanyiBatman/tagSonic/internal/scanner"
	"github.com/joho/godotenv"
)

type Song struct {
	FilePath          string
	FingerprintResult acoustid.FingerprintResult
	OldTags           id3tagger.Tags
	SongMetadata      acoustid.SongMetadata
	ConfidenceScore   int
}

// TaggedSong holds information about a successfully tagged song.
type TaggedSong struct {
	OriginalFilename string
	NewTitle         string
	NewArtist        string
	ConfidenceScore  int
}

// FailedSong holds information about a song that could not be tagged.
type FailedSong struct {
	OriginalFilename string
	Reason           string
}

const confidenceThreshold = 85

// calculateConfidence scores how confident we are that the fetched tags
// actually belong to this file. It checks the AcoustID score, old ID3 tags,
// and the raw filename to build up a total confidence score.
func calculateConfidence(songInfo Song, newTags *id3tagger.Tags) int {
	// 1. Base score from AcoustID (0-100)
	score := int(songInfo.SongMetadata.Score * 100)

	// 2. Bonus if old ID3 title matches the new title
	if strings.EqualFold(songInfo.OldTags.Title, newTags.Title) {
		score += 10
	}

	// 3. Bonus if old ID3 artist matches the new artist
	if strings.EqualFold(songInfo.OldTags.Artist, newTags.Artist) {
		score += 10
	}

	// 4. Bonus if the raw filename contains the new title
	rawFilename := filepath.Base(songInfo.FilePath)
	rawFilename = strings.TrimSuffix(rawFilename, filepath.Ext(rawFilename))

	lowFilename := strings.ToLower(rawFilename)
	lowNewTitle := strings.ToLower(newTags.Title)

	// Add massive bonus for filename match, since that's a very strong human signal
	if strings.Contains(lowFilename, lowNewTitle) || strings.Contains(lowNewTitle, lowFilename) {
		score += 30
	} else if strings.EqualFold(songInfo.OldTags.Title, "") && strings.EqualFold(songInfo.OldTags.Artist, "") {
		// If there are no old tags AND the filename doesn't match the new title at all, penalize
		score -= 20
	}

	// 5. Hard Title Penalty: If the new Title is NOT in the filename AND NOT in the old ID3 title
	// This is the most common cause of false positives (correct artist, wrong song).
	if !strings.EqualFold(songInfo.OldTags.Title, newTags.Title) &&
		!strings.Contains(lowFilename, lowNewTitle) &&
		!strings.Contains(lowNewTitle, lowFilename) {
		score -= 40
	}

	// 6. Active Penalty: If literally nothing matches (no artist match, no title match, no filename match)
	if !strings.EqualFold(songInfo.OldTags.Title, newTags.Title) &&
		!strings.EqualFold(songInfo.OldTags.Artist, newTags.Artist) &&
		!strings.Contains(lowFilename, lowNewTitle) {
		score -= 30
	}

	return score
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	acousticIDAPIKey := os.Getenv("AcousticIDAPIKey")
	fmt.Println("Hello World")
	targetDir := "./audios"
	if len(os.Args) > 1 {
		targetDir = os.Args[1]
	}

	fmt.Printf("Scanning directory: %s\n", targetDir)
	songs, err := scanner.ScanDir(targetDir)
	if err != nil {
		fmt.Println(err)
	}

	var taggedSongs []TaggedSong
	var failedSongs []FailedSong

	for _, song := range songs {
		fmt.Printf("Processing: %s\n", song)

		fp, err := acoustid.Fingerprint(song)
		if err != nil {
			fmt.Printf("  ✗ Fingerprint failed: %v\n\n", err)
			failedSongs = append(failedSongs, FailedSong{OriginalFilename: song, Reason: "Fingerprint failed"})
			continue
		}
		tags, err := id3tagger.GetID3Tags(song)
		if err != nil {
			fmt.Printf("  ✗ Could not read existing tags: %v\n\n", err)
			failedSongs = append(failedSongs, FailedSong{OriginalFilename: song, Reason: "Could not read existing tags"})
			continue
		}
		songInfo := Song{
			FilePath:          song,
			FingerprintResult: fp,
			OldTags:           *tags,
		}
		lookupData, err := acoustid.LookupMetadata(fp.Fingerprint, fp.Duration, acousticIDAPIKey)
		if err != nil {
			fmt.Printf("  ✗ AcoustID lookup failed: %v\n\n", err)
			failedSongs = append(failedSongs, FailedSong{OriginalFilename: song, Reason: "AcoustID lookup failed"})
			continue
		}
		songInfo.SongMetadata = *lookupData
		fmt.Printf("  AcoustID matched: %s - %s (score: %.0f%%)\n", lookupData.Title, lookupData.Artist, lookupData.Score*100)

		// Fetch tags from iTunes first
		var newTags *id3tagger.Tags
		itunesTags, err := itunes.LookupSong(songInfo.SongMetadata.Title, songInfo.SongMetadata.Artist)
		if err == nil {
			newTags = &id3tagger.Tags{
				Title:       itunesTags.Title,
				Artist:      itunesTags.Artist,
				Album:       itunesTags.Album,
				Genre:       itunesTags.Genre,
				Year:        itunesTags.Year,
				CoverArtURL: itunesTags.CoverArtURL,
			}
			fmt.Printf("  ✓ Fetched tags from iTunes\n")
		} else {
			// Fallback to MusicBrainz
			mbTags, err := musicbrainz.LookupByMBID(songInfo.SongMetadata.RecordingMBID)
			if err == nil {
				newTags = &id3tagger.Tags{
					Title:       mbTags.Title,
					Artist:      mbTags.Artist,
					Album:       mbTags.Album,
					Genre:       mbTags.Genre,
					Year:        mbTags.Year,
					CoverArtURL: mbTags.CoverArtURL,
				}
				fmt.Printf("  ✓ Fetched tags from MusicBrainz\n")
			} else {
				fmt.Printf("  ✗ Could not fetch tags from any source\n")
			}
		}

		if newTags == nil {
			failedSongs = append(failedSongs, FailedSong{OriginalFilename: song, Reason: "Could not fetch tags from iTunes or MusicBrainz"})
			continue
		}

		// Calculate confidence
		songInfo.ConfidenceScore = calculateConfidence(songInfo, newTags)

		// Fallback: If AcoustID matched something totally wrong (low confidence),
		// try a direct text search using the filename + old tags instead of the audio fingerprint.
		if songInfo.ConfidenceScore < confidenceThreshold {
			fmt.Printf("  ⚠ AcoustID match confidence too low (%d%%). Trying pure text fallback search...\n", songInfo.ConfidenceScore)

			rawFilename := filepath.Base(songInfo.FilePath)
			rawFilename = strings.TrimSuffix(rawFilename, filepath.Ext(rawFilename))

			// Build a robust text query prioritizing existing tags, fallback to filename
			searchQuery := strings.TrimSpace(songInfo.OldTags.Title + " " + songInfo.OldTags.Artist)
			if len(searchQuery) < 3 {
				searchQuery = rawFilename
			}

			fallbackTags, err := itunes.LookupSong(searchQuery, "")
			if err == nil {
				newTags = &id3tagger.Tags{
					Title:       fallbackTags.Title,
					Artist:      fallbackTags.Artist,
					Album:       fallbackTags.Album,
					Genre:       fallbackTags.Genre,
					Year:        fallbackTags.Year,
					CoverArtURL: fallbackTags.CoverArtURL,
				}

				// Re-score based on the text search result. We don't have a base score from AcoustID anymore,
				// so we start at 50 and rely heavily on the text/filename matching bonuses.
				songInfo.SongMetadata.Score = 0.50
				songInfo.ConfidenceScore = calculateConfidence(songInfo, newTags)
				fmt.Printf("  ✓ Text fallback found: %s - %s\n", newTags.Title, newTags.Artist)
			} else {
				fmt.Printf("  ✗ Text fallback search failed: %v\n", err)
			}
		}

		fmt.Printf("  Final Confidence: %d%% (threshold: %d%%)\n", songInfo.ConfidenceScore, confidenceThreshold)

		if songInfo.ConfidenceScore < confidenceThreshold {
			fmt.Printf("  ⚠ SKIPPED — confidence too low, not overwriting tags\n\n")
			failedSongs = append(failedSongs, FailedSong{
				OriginalFilename: song,
				Reason:           fmt.Sprintf("Confidence too low (%d%% < %d%%)", songInfo.ConfidenceScore, confidenceThreshold),
			})
			continue
		}

		err = id3tagger.WriteID3Tags(newTags, song)
		if err != nil {
			fmt.Printf("  ✗ Error writing tags: %v\n\n", err)
			failedSongs = append(failedSongs, FailedSong{OriginalFilename: song, Reason: fmt.Sprintf("Failed to write tags: %v", err)})
		} else {
			fmt.Printf("  ✓ Successfully updated tags\n\n")
			taggedSongs = append(taggedSongs, TaggedSong{
				OriginalFilename: song,
				NewTitle:         newTags.Title,
				NewArtist:        newTags.Artist,
				ConfidenceScore:  songInfo.ConfidenceScore,
			})
		}
	}

	// Write log report
	writeLogReport(taggedSongs, failedSongs)
}

func writeLogReport(tagged []TaggedSong, failed []FailedSong) {
	file, err := os.Create("tagSonic_log.txt")
	if err != nil {
		fmt.Printf("Failed to create log report: %v\n", err)
		return
	}
	defer file.Close()

	file.WriteString("=== tagSonic Execution Log ===\n\n")

	file.WriteString(fmt.Sprintf("Successfully Tagged: %d songs\n", len(tagged)))
	file.WriteString("----------------------------------------\n")
	for _, song := range tagged {
		file.WriteString(fmt.Sprintf("[%d%%] %s -> %s - %s\n", song.ConfidenceScore, song.OriginalFilename, song.NewArtist, song.NewTitle))
	}

	file.WriteString(fmt.Sprintf("\nFailed or Skipped: %d songs\n", len(failed)))
	file.WriteString("----------------------------------------\n")
	for _, song := range failed {
		file.WriteString(fmt.Sprintf("[ERROR] %s (Reason: %s)\n", song.OriginalFilename, song.Reason))
	}

	fmt.Printf("Execution finished. Wrote report to tagSonic_log.txt\n")
}
