package util

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	uuid "github.com/satori/go.uuid"
)

var (
	started, traps, exiting bool
)

func IsSocket(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.Mode().Type() == fs.ModeSocket
}

func HexToInt(hexString string) int {

	if hexString[:2] != "0x" {
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

func setTraps() {
	if !traps {
		traps = true
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		go func() {
			c := <-sigc
			DispatchEvent(SignalEvent, c.String())
			Shutdown(0)
		}()
	}
}

// SetTraps to be run first in your application
func Initialize() {
	if !started {
		started = true
		log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
		flag.Parse()
		setTraps()
		DispatchEvent(InitializeEvent)
	}
}

func Shutdown(code int) {
	if !exiting {
		exiting = true
		DispatchEvent(ShutdownEvent, code)
		os.Exit(code)
	}
}
