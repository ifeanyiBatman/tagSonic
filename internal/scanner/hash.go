package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func HashFiles(filePaths []string) ([]string, error) {
	var hashes []string
	for _, filePath := range filePaths {
		hash, err := hashFile(filePath)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, hash)
	}
	for _, hash := range hashes {
		fmt.Println(hash)
	}
	return hashes, nil
}

func hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
