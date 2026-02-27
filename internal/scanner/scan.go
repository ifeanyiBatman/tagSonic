package scanner

import (
	"fmt"
	"os"
	"path/filepath"
)

func ScanDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			if IsAudioFile(entry) {
				files = append(files, filepath.Join(dir, entry.Name()))
			}
		} else {
			subDirFiles, err := ScanDir(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, subDirFiles...)
		}
	}
	fmt.Printf("You have %d songs\n", len(files))
	for _, file := range files {
		fmt.Println(file)
	}
	return files, nil
}
