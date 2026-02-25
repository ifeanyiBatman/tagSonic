package readFolder

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadDirTolist(dir string) ([]string, error) {

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			if FileISAudio(entry) {
				files = append(files, filepath.Join(dir, entry.Name()))
			}
		} else {
			subDirFiles, err := ReadDirTolist(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, subDirFiles...)
		}
	}
	fmt.Printf("You have %d songs\n",len(files))
	for _, file := range files {
		fmt.Println(file)
	}
	return files, nil
}