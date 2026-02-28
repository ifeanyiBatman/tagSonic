package scanner

import (
	"os"
	"path/filepath"
)

func IsAudioFile(file os.DirEntry) bool {
	filename := file.Name()
	fileExtension := filepath.Ext(filename)
	switch fileExtension {
//  case ".mp3", ".wav", ".flac", ".ogg", ".m4a": we only support mp3 for now
	case ".mp3":
		return true
	default:
		return false
	}
}
