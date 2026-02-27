package main

import (
	"fmt"

	"github.com/ifeanyiBatman/tagSonic/internal/fingerprinting"
	"github.com/ifeanyiBatman/tagSonic/internal/readFolder"
)

func main() {
	fmt.Println("Hello World")
	songs, err := readFolder.ReadDirTolist("./audios")
	if err != nil {
		fmt.Println(err)
	}
	readFolder.HashFiles(songs)
	fp, err := fingerprinting.Fingerprint("audios/kanyewest/30 Hours.mp3")
	fmt.Printf("fingerprint for the track 30 hours\n%s", fp.FingerPrint)
	id, err := fingerprinting.GetIdFromFingerprint(fp.FingerPrint, fp.Duration, "JYBesruRxAQ")
	if err != nil {
		fmt.Printf("Error looking up fingerprint: %v\n", err)
		return
	}
	fmt.Printf("\nMatched AcoustID for the track 30 hours: %s\n", id)
}
