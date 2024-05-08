package util

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"

	uuid "github.com/satori/go.uuid"
)

func IsSocket(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.Mode().Type() == fs.ModeSocket
}

func HexToInt(hexString string) int {

	if fmt.Sprintf("%s  ", hexString)[:2] != "0x" {
		hexString = "0x" + hexString
	}

	i, err := strconv.ParseInt(hexString, 0, 64)
	if err != nil {
		return 0
	}
	return int(i)
}

func GenerateUid() string {
	return uuid.NewV4().String()
}
