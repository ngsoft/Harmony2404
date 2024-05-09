package main

import (
	"flirc/usocket"
	"flirc/util"
	"flirc/wsocket"
)

func main() {
	util.Initialize()

	logger.Info(
		"flags(port=>%v, socket=>%v, remote => %v, delay => %v, ping => %v)",
		*wsPort,
		*socket,
		*remote,
		*keyDelay,
		*pingOn,
	)

	flirc = FlircHandler{
		remote:  *remote,
		delay:   *keyDelay,
		keymaps: LoadKeymaps(),
	}
	flirc.SetLoggerPrefix("[usocket]")
	s, ok := usocket.ConnectSocket(*socket, &flirc)
	if ok {
		flirc.Info("connected to %s", *socket)
		util.AddEventHandler(&flirc)
		go s.Run()
	}

	// launch WS
	ws.Config = wsocket.Config{
		Ping: *pingOn,
		Port: *wsPort,
	}
	flirc.Room = ws.AddRoom("remote")

	ws.Run()
}
