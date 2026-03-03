package id3tagger

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bogem/id3v2/v2"
)

func WriteID3Tags(newTags *Tags, filepath string) error {
	tags, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer tags.Close()

	// Use UTF-8 for everything to avoid "rune not supported" errors
	tags.SetDefaultEncoding(id3v2.EncodingUTF8)

	tags.SetTitle(newTags.Title)
	tags.SetArtist(newTags.Artist)
	tags.SetAlbum(newTags.Album)
	tags.SetGenre(newTags.Genre)
	tags.SetYear(newTags.Year)

	// Fetch and attach cover art if URL is present
	if newTags.CoverArtURL != "" {
		resp, err := http.Get(newTags.CoverArtURL)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch cover art from %s: %v\n", newTags.CoverArtURL, err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				imgBytes, err := io.ReadAll(resp.Body)
				if err == nil {
					mimeType := "image/jpeg"
					if resp.Header.Get("Content-Type") != "" {
						mimeType = resp.Header.Get("Content-Type")
					}

					pic := id3v2.PictureFrame{
						Encoding:    id3v2.EncodingUTF8,
						MimeType:    mimeType,
						PictureType: id3v2.PTFrontCover,
						Description: "Front cover",
						Picture:     imgBytes,
					}
					tags.AddAttachedPicture(pic)
				}
			} else {
				fmt.Printf("Warning: Cover art URL returned status %d\n", resp.StatusCode)
			}
		}
	}

	err = tags.Save()
	if err != nil {
		return err
	}
	return nil
}
