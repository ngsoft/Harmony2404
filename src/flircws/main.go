package main

import (
	"flirc/usocket"
	"flirc/util"
	"log"
)

func main() {
	util.Initialize()

	log.Printf(
		"INFO: flags(port=>%v, socket=>%v, remote => %v, delay => %v)",
		*wsPort,
		*socket,
		*remote,
		*keyDelay,
	)

	flirc = FlircHandler{
		remote:  *remote,
		delay:   *keyDelay,
		keymaps: LoadKeymaps(),
	}
	s, ok := usocket.ConnectSocket(*socket, &flirc)
	if ok {
		logger.Info("connected to %s", *socket)
		util.AddEventHandler(&flirc)
		s.Run()

	}
	for s.HasStatus() {
		continue
	}
}
