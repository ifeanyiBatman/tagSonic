package main

import (
	"fmt"
	"os"

	"github.com/ifeanyiBatman/tagSonic/internal/acoustid"
	"github.com/ifeanyiBatman/tagSonic/internal/id3tagger"
	"github.com/ifeanyiBatman/tagSonic/internal/scanner"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	acousticIDAPIKey := os.Getenv("AcousticIDAPIKey")
	fmt.Println("Hello World")
	songs, err := scanner.ScanDir("./audios")
	if err != nil {
		fmt.Println(err)
	}
	scanner.HashFiles(songs)
	fp, err := acoustid.Fingerprint("audios/kanyewest/30 Hours.mp3")
	fmt.Printf("fingerprint for the track 30 hours\n%s", fp.Fingerprint)
	meta, err := acoustid.LookupMetadata(fp.Fingerprint, fp.Duration, acousticIDAPIKey)
	if err != nil {
		fmt.Printf("Error looking up fingerprint: %v\n", err)
		return
	}
	fmt.Printf("\nMatched AcoustID for the track 30 hours: %s\n", meta.ID)
	fmt.Printf("Title: %s\n", meta.Title)
	fmt.Printf("Artist: %s\n", meta.Artist)

	tags, err := id3tagger.GetID3Tags("audios/kanyewest/30 Hours.mp3")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Title: %s\n", tags.Title)
	fmt.Printf("Artist: %s\n", tags.Artist)
	fmt.Printf("Album: %s\n", tags.Album)
	fmt.Printf("Genre: %s\n", tags.Genre)
	fmt.Printf("Year: %s\n", tags.Year)
}
