package main

import (
	"flag"
	"log"
)

const (
	defaultPort   int    = 9023
	defaultSocket string = "/var/run/lirc/lircd"
)

var (
	socket *string
	wsPort *int
)

func parseFlags() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	wsPort = flag.Int("port", defaultPort, "Websocket listen port")
	socket = flag.String("socket", defaultSocket, "InputLirc unix socket location")
	flag.Parse()

	log.Printf(
		"INFO: flags(port=>%v, socket=>`%v`)",
		*wsPort,
		*socket,
	)
}
