package main

import (
	"flirc/usocket"
	"flirc/util"
	"flirc/wsocket"
	"time"
)

func main() {

	// load the flags, the traps, app events
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
		remote:     *remote,
		delay:      *keyDelay,
		keymaps:    LoadKeymaps(),
		UnixSocket: &usocket.UnixSocket{},
	}
	flirc.SetLoggerPrefix("[usocket]")
	// launch WS
	ws.Config = wsocket.Config{
		Ping:  *pingOn,
		Port:  *wsPort,
		Route: wsRoute,
	}
	flirc.Room = ws.AddRoom("remote")
	ws.AddHandler(&flirc)
	// if socket is disconnected when running app

	go func() {
		var (
			ticker = time.NewTicker(5 * time.Second)
		)
		defer ticker.Stop()
		for ; true; <-ticker.C {
			if util.IsStopping() {
				return
			}
			s, ok := usocket.ConnectSocket(*socket, &flirc)
			if ok {
				flirc.Info("connected to %s", *socket)
				go s.Run()
				return
			}

		}
	}()

	// to use login middleware
	// users := wsocket.InMemoryUserValidator{}
	// users.AddUsers(wsocket.NewUser("admin", "pwd"))
	// ws.AddHandler(wsocket.NewLoginHandler(&ws, 20, &users))
	ws.Run() // loading in the main thread to nt kill the app
}
