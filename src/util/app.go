package util

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	globalEvents            EventListener
	started, traps, exiting bool
)

func IsStarting() bool {
	return !started
}
func IsRunning() bool {
	return started && !exiting
}
func IsStopping() bool {
	return exiting
}

func AddEventHandler(e EventHandler, v ...EventType) func() {
	return globalEvents.AddEventHandler(e, v...)
}
func RemoveEventHandler(e EventHandler, v ...EventType) {
	globalEvents.RemoveEventHandler(e, v...)
}
func DispatchEvent(e EventType, p ...interface{}) {
	globalEvents.DispatchEvent(e, p...)
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

// Initialize to be run first in your application
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
