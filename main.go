package main

import (
	"fmt"
	"os"

	"github.com/ifeanyiBatman/tagSonic/internal/fingerprinting"
	"github.com/ifeanyiBatman/tagSonic/internal/readFolder"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	acousticIDAPIKey := os.Getenv("AcousticIDAPIKey")
	fmt.Println("Hello World")
	songs, err := readFolder.ReadDirTolist("./audios")
	if err != nil {
		fmt.Println(err)
	}
	readFolder.HashFiles(songs)
	fp, err := fingerprinting.Fingerprint("audios/kanyewest/30 Hours.mp3")
	fmt.Printf("fingerprint for the track 30 hours\n%s", fp.FingerPrint)
	meta, err := fingerprinting.GetMetadata(fp.FingerPrint, fp.Duration, acousticIDAPIKey)
	if err != nil {
		fmt.Printf("Error looking up fingerprint: %v\n", err)
		return
	}
	fmt.Printf("\nMatched AcoustID for the track 30 hours: %s\n", meta.ID)
	fmt.Printf("Title: %s\n", meta.Title)
	fmt.Printf("Artist: %s\n", meta.Artist)
}
