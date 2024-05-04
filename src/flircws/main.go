package main

import (
	"flirc/util"
	"fmt"
	"log"
	"net/http"
)

func main() {
	util.SetTraps()
	util.ParseFlags()
	log.Printf(
		"INFO: flags(port=>%v, socket=>%v, remote => %v)",
		*wsPort,
		*socket,
		*remote,
	)

	flirc = FlircHandler{
		remote: *remote,
		path:   *socket,
		BaseHandler: util.BaseHandler{
			Uid:    util.GenerateUid(),
			Logger: util.NewLogger(fmt.Sprintf("[%s]", *remote)),
		},
	}

	flirc.AddEventHandler(&flirc, INPUT_EVENT)

	// connect flirc
	s, ok := util.NewUnixSocket(*socket, &flirc)
	if !ok {
		logger.Error("cannot connect to socket %s", *socket)
		util.Shutdown(1)
	}

	go s.HandleConnection()

	// connect webhook

	ws = ConnHandler{
		Port:     *wsPort,
		Route:    wsRoute,
		handlers: make(map[string]Handler, 0),
	}

	util.CreateWebSocketRoute(ws.Route, &ws)

	logger.Info("Listening to tcp port %v %v", ws.Port, ws.Route)
	if err := http.ListenAndServe(ws.GetHttpPort(), nil); err != nil {
		logger.Error(err.Error())
	}

}
