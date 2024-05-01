package main

const (
	defaultPort   int    = 9023
	defaultSocket string = "/var/run/lirc/lircd"
	defaultRemote string = "FLIRC"
	WsRoute       string = "/ws"
)

var (
	socket, remote *string
	wsPort         *int
)
