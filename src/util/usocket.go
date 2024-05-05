package util

import (
	"fmt"
	"net"
	"strings"
)

type UnixSocketHandler interface {
	HandleUnixSocket(*UnixSocket)
}

type UnixSocket struct {
	Status
	Path    string
	conn    net.Conn
	handler UnixSocketHandler
	ReadWriteLock
	EventListener
	Logger
}

func NewUnixSocket(path string, handler UnixSocketHandler) (*UnixSocket, bool) {

	var socket UnixSocket

	if !IsSocket(path) {
		return &socket, false
	}

	c, err := net.Dial("unix", path)
	if err != nil {
		return &socket, false
	}

	socket = UnixSocket{
		Status:  STATUS_ON,
		conn:    c,
		Path:    path,
		handler: handler,
		Logger:  NewLogger(fmt.Sprintf("[%s]", path)),
	}

	AddEventHandler(&socket, ShutdownEvent)

	if s, ok := handler.(EventHandler); ok {
		socket.AddEventHandler(s)
	}

	return &socket, true
}
func (h *UnixSocket) OnEvent(ev *Event) {
	switch ev.Type {
	case ShutdownEvent:
		h.Info("received shutdown event, closing connection")
		h.Close()
	}
}

func (h *UnixSocket) Close() {

	if h.On() {
		h.Status = STATUS_OFF
		h.Info("Closing Connection")
		_ = h.conn.Close()
		h.Info("Connection closed")
		h.DispatchEvent(CONNECTION_CLOSE, h.Path)
	}

}

func (h *UnixSocket) HandleConnection() {
	if h.On() {
		defer h.Close()
		h.handler.HandleUnixSocket(h)
	}

}

func (h *UnixSocket) ReadMessage() (string, bool) {
	var (
		l     int
		msg   string
		input = make([]byte, 1024)
		err   = fmt.Errorf("")
	)

	if h.On() {
		h.ReadLock.Lock()
		defer h.ReadLock.Unlock()

		l, err = h.conn.Read(input)
		if err != nil || l == 0 {
			h.Close()
			return msg, false
		}
		msg = strings.TrimRight(string(input[:l]), "\r\n")
	}

	return msg, err == nil
}

func (h *UnixSocket) WriteMessage(message string) bool {

	if len(message) == 0 {
		return false
	}

	var (
		l   int
		err error
		msg []byte
	)

	if h.On() {
		h.WriteLock.Lock()
		defer h.WriteLock.Unlock()

		msg = []byte(strings.TrimRight(message, "\r\n") + "\n")

		l, err = h.conn.Write(msg)
		if err != nil {
			h.Close()
			return false
		}
		return l > 0
	}
	return false
}
