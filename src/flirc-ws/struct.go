package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

type BaseHandler struct {
	Uid    string
	Logger *log.Logger
	Mutex  struct {
		Read  sync.Mutex
		Write sync.Mutex
	}
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
	*BaseHandler
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

func (h *BaseHandler) log(format string, v ...interface{}) {
	if h.Logger == nil {
		h.Logger = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	}

	pfx := h.Uid[:10]
	if pfx != "" {
		pfx = "[" + pfx + "] "
	}
	h.Logger.Output(2, fmt.Sprintf(pfx+format, v...))
}

func (h *BaseHandler) error(format string, v ...interface{}) {
	h.log("ERR: "+format, v...)
}
func (h *BaseHandler) info(format string, v ...interface{}) {
	h.log("INFO: "+format, v...)
}
