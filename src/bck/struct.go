package main

import (
	"flirc/util"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

type SocketMessage struct {
	Message string
	Code    string
	Repeat  int
	Key     string
	Remote  string
}

type UnixSocket struct {
	Conn net.Conn
	Open bool
}

type WebSocket struct {
	Upgrader websocket.Upgrader
	Writer   http.ResponseWriter
	Conn     *websocket.Conn
	Open     bool
}

type Handler struct {
	*util.BaseHandler
	Remote     string
	Socket     string
	Port       int
	Route      string
	UnixSocket UnixSocket
	WebSocket  WebSocket
}

func newHandler() *Handler {

	var handler Handler = Handler{
		BaseHandler: &BaseHandler{
			Uid: generateUid(),
		},
		Socket: *socket,
		Port:   *wsPort,
		Remote: *remote,
		Route:  WsRoute,
		WebSocket: WebSocket{
			Upgrader: websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin:     func(r *http.Request) bool { return true },
			},
		},
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		handler.Shutdown(0)
	}()

	return &handler
}
