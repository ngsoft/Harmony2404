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
	GlobalEvents EventListener
	traps        bool
	exiting      bool
)

func IsSocket(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.Mode().Type() == fs.ModeSocket
}

func HexToInt(hexString string) int {
	i, err := strconv.ParseInt("0x"+hexString, 0, 64)
	if err != nil {
		return 0
	}
	return int(i)
}

func GenerateUid() string {
	return uuid.NewV4().String()
}

func ParseFlags() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
}

// SetTraps to be run first in your application
func SetTraps() {
	if !traps {
		traps = true
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		go func() {
			c := <-sigc
			GlobalEvents.DispatchEvent(SIGNAL_EVENT, c.String())
			Shutdown(0)
		}()
	}
}

func Shutdown(code int) {
	if !exiting {
		exiting = true
		GlobalEvents.DispatchEvent(SHUTDOWN_EVENT, strconv.Itoa(code))
		os.Exit(code)
	}

}
