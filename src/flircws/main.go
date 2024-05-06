package main

import (
	"flirc/util"
	"log"
)

func main() {
	util.Initialize()

	log.Printf(
		"INFO: flags(port=>%v, socket=>%v, remote => %v)",
		*wsPort,
		*socket,
		*remote,
	)

	flirc = FlircHandler{
		remote: *remote,
		path:   *socket,
	}

	flirc.Initialize()

	// flirc.AddEventHandler(&flirc, INPUT_EVENT)

	// connect flirc
	s, ok := util.NewUnixSocket(*socket, &flirc)
	if !ok {
		logger.Error("cannot connect to socket %s", *socket)
		util.Shutdown(1)
	}

	s.HandleConnection()

	// connect webhook

	// ws = ConnHandler{
	// 	Port:     *wsPort,
	// 	Route:    wsRoute,
	// 	handlers: make(map[string]Handler, 0),
	// }

	// util.CreateWebSocketRoute(ws.Route, &ws)

	// logger.Info("Listening to tcp port %v %v", ws.Port, ws.Route)
	// if err := http.ListenAndServe(ws.GetHttpPort(), nil); err != nil {
	// 	logger.Error(err.Error())
	// }

}
