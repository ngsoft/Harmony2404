package usocket

import (
	"bytes"
	"flirc/util"
	"net"
	"strings"
	"time"
)

const (
	Incoming util.EventType = "usocket.in"
	Outgoing util.EventType = "usocket.out"
	Open     util.EventType = "usocket.open"
	Close    util.EventType = "usocket.close"
)

var (
	newLine = "\n"
)

type SocketEvent struct {
	Type    util.EventType
	Message string
}

type Handler interface {
	OnSocketEvent(*SocketEvent, *UnixSocket)
}

type UnixSocket struct {
	util.Status
	Handler
	util.ReadWriteLock
	init bool
	file string
	conn net.Conn
}

func ConnectSocket(file string, h Handler) (*UnixSocket, bool) {
	var socket UnixSocket
	if !util.IsSocket(file) {
		return &socket, false
	}
	c, err := net.Dial("unix", file)
	if err != nil {
		return &socket, false
	}
	socket = UnixSocket{
		Handler: h,
		file:    file,
		conn:    c,
		Status:  util.NewStatus(),
	}
	return &socket, true
}

func (s *UnixSocket) Run() {
	if s.file == "" {
		panic("UnixSocket was not instantiated using ConnectSocket")
	}
	if s.init {
		return
	}
	if s.On() {
		go s.OnSocketEvent(&SocketEvent{Type: Open}, s)
		s.init = true
		go s.read()
	}

}

func (s *UnixSocket) Close() {
	if s.On() {
		s.Status++
		_ = s.conn.Close()
		s.OnSocketEvent(&SocketEvent{Type: Close}, s)
	}

}

func (s *UnixSocket) read() {
	var (
		l     int
		input = make([]byte, 1024)
		err   error
	)

	s.ReadLock.Lock()
	defer s.ReadLock.Unlock()

	for s.On() {
		l, err = s.conn.Read(input)
		if err != nil || l == 0 {
			s.Close()
			return
		}
		for _, line := range bytes.Split(input[:l], []byte(newLine)) {
			if len(line) > 0 {
				go s.OnSocketEvent(&SocketEvent{
					Type:    Incoming,
					Message: string(line),
				}, s)
			}

		}
	}
}

func (s *UnixSocket) SendMessage(message string) bool {

	if s.On() && len(message) > 0 {
		s.WriteLock.Lock()
		defer s.WriteLock.Unlock()
		var (
			l   int
			err error
		)
		for _, line := range strings.Split(message, newLine) {
			if len(line) > 0 {
				l, err = s.conn.Write([]byte(line + newLine))
				if err != nil || l == 0 {
					s.Close()
					return false
				}
				go s.OnSocketEvent(&SocketEvent{
					Type:    Outgoing,
					Message: line,
				}, s)
			}

		}
		return true
	}
	return false
}

func (s *UnixSocket) Reconnect(delay time.Duration) {
	if !s.Off() || !s.init {
		return
	}

	var (
		ticker = time.NewTicker(delay)
	)

	defer ticker.Stop()

	for range ticker.C {
		if !util.IsStopping() {
			c, err := net.Dial("unix", s.file)
			if err == nil {
				s.conn = c
				s.Status = util.NewStatus()
				go s.OnSocketEvent(&SocketEvent{Type: Open}, s)
				go s.read()
				return
			}
		}

	}
}
