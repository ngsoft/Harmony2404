package main

import (
	"flag"
	"log"
)

func parseFlags() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	wsPort = flag.Int("port", defaultPort, "Websocket listen port")
	socket = flag.String("socket", defaultSocket, "InputLirc unix socket location")
	remote = flag.String("remote", defaultRemote, "InputLirc Remote channel")
	flag.Parse()
	log.Printf(
		"INFO: flags(port=>%v, socket=>%v, remote => %v)",
		*wsPort,
		*socket,
		*remote,
	)
}
