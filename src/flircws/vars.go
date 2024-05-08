package main

import (
	"flag"
	"flirc/util"
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
	flirc    FlircHandler
	logger   = util.NewLogger("[MAIN]")
	etc      = []string{
		"../../etc",
		"../../../etc",
	}
	cfgDir  = "flircd"
	keymaps []Keymap
)
