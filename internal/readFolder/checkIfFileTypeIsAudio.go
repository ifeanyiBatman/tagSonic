package readFolder

import (
	"os"
	"path/filepath"
)

func FileISAudio(file os.DirEntry) bool {
	filename := file.Name()
	fileExtension := filepath.Ext(filename)
	switch fileExtension {
		case ".mp3", ".wav", ".flac", ".ogg", ".m4a":
			return true
		default:
			return false
	}
}