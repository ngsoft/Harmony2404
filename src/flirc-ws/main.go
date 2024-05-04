package main

import (
	"flirc/util"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	util.ParseFlags()
	log.Printf(
		"INFO: flags(port=>%v, socket=>%v, remote => %v)",
		*wsPort,
		*socket,
		*remote,
	)

	// connect flirc
	c, ok := util.NewUnixSocket(*socket)
	if !ok {
		log.Fatal("cannot connect to socket ", *socket)
		os.Exit(1)
	}

	flirc = FlircHandler{
		remote:      *remote,
		UnixSocket:  c,
		BaseHandler: util.NewBaseHandler(),
	}
	flirc.AttachEvent(func(e *util.Event) {
		flirc.Log("%v", flirc.Message)
	}, INPUT_EVENT)

	go c.HandleConnection(flirc.handleConn)

	// connect webhook

	route, _ := util.NewRouteHandler(WsRoute)
	// route.SetRequestHandler(HandleWebookRoute)

	flirc.Log("Listening to tcp port %v %v", *wsPort, route)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *wsPort), nil); err != nil {
		flirc.Error(err.Error())
	}

}
