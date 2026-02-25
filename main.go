package main

import (
	"fmt"
	"github.com/ifeanyiBatman/tagSonic/internal/readFolder"
)

func main() {
	fmt.Println("Hello World")
	readFolder.ReadDirTolist("./audios")
}