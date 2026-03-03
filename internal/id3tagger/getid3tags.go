package id3tagger

import (
	"fmt"

	"github.com/bogem/id3v2/v2"
)

type Tags struct {
	Title       string
	Artist      string
	Album       string
	Genre       string
	Year        string
	CoverArtURL string
}

func GetID3Tags(filepath string) (*Tags, error) {
	tags, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer tags.Close()
	return &Tags{
		Title:  tags.Title(),
		Artist: tags.Artist(),
		Album:  tags.Album(),
		Genre:  tags.Genre(),
		Year:   tags.Year(),
	}, nil
}
