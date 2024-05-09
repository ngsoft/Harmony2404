package main

import (
	"flag"
	"flirc/util"
	"flirc/wsocket"
)

const (
	defaultPort   int    = 9023
	defaultSocket string = "/var/run/lirc/lircd"
	defaultRemote string = "FLIRC"
	defaultDelay  int    = 400
	wsRoute       string = "/ws"
)

var (
	wsPort   *int    = flag.Int("port", defaultPort, "Websocket listen port")
	socket   *string = flag.String("socket", defaultSocket, "InputLirc unix socket location")
	remote   *string = flag.String("remote", defaultRemote, "InputLirc Remote channel")
	keyDelay *int    = flag.Int("delay", defaultDelay, "Inter key delay in ms")
	pingOn   *bool   = flag.Bool("ping", false, "Enable web socket ping")
	flirc    FlircHandler
	logger   util.Logger
	cfgDir   = "etc/flircd"
	// libDir   = "usr/local/lib/flircd"
	ws wsocket.WebSocket
)
