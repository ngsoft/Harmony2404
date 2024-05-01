package main

import (
	"io/fs"
	"log"
	"os"
)

func IsSocket(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal("ERR: ", err)
	}
	return fileInfo.Mode().Type() == fs.ModeSocket
}
