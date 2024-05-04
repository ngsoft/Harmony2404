package main

import (
	"flag"
	"flirc/util"
)

const (
	defaultPort   int    = 9023
	defaultSocket string = "/var/run/lirc/lircd"
	defaultRemote string = "FLIRC"
	wsRoute       string = "/ws"
)

var (
	wsPort *int    = flag.Int("port", defaultPort, "Websocket listen port")
	socket *string = flag.String("socket", defaultSocket, "InputLirc unix socket location")
	remote *string = flag.String("remote", defaultRemote, "InputLirc Remote channel")
	flirc  FlircHandler
	ws     ConnHandler
	logger = util.NewLogger("[MAIN]")
)
