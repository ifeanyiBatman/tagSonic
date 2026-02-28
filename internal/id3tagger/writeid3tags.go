package id3tagger

import (
	"fmt"
	"github.com/bogem/id3v2/v2"
)

func WriteID3Tags(newTags *Tags, filepath string) error {
	tags, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer tags.Close()
	tags.SetTitle(newTags.Title)
	tags.SetArtist(newTags.Artist)
	tags.SetAlbum(newTags.Album)
	tags.SetGenre(newTags.Genre)
	tags.SetYear(newTags.Year)
	return nil
}
