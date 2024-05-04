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
	Status  Status
	Path    string
	conn    net.Conn
	handler UnixSocketHandler
	ReadWriteLock
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
	}
	GlobalEvents.AddEventHandler(&socket, SHUTDOWN_EVENT)
	return &socket, true
}
func (h *UnixSocket) HandleEvent(ev *Event) {
	switch ev.Type {
	case SHUTDOWN_EVENT:
		h.Close()
	}
}

func (h *UnixSocket) Close() {
	if h.Status.On() {
		_ = h.conn.Close()
	}
	h.Status = STATUS_OFF
}

func (h *UnixSocket) HandleConnection(handler UnixSocketHandler) {
	if h.Status.On() {
		defer h.Close()
		handler.HandleUnixSocket(h)
	}

}

func (h *UnixSocket) ReadMessage() (string, bool) {
	var (
		l     int
		msg   string
		input []byte = make([]byte, 1024)
		err   error  = fmt.Errorf("")
	)

	if h.Status.On() {
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

	if h.Status.On() {
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
