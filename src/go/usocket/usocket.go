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

type UnixSocket struct {
	util.Status
	util.EventListener
	util.ReadWriteLock
	init bool
	file string
	conn net.Conn
}

func ConnectSocket(file string, h util.EventHandler) (*UnixSocket, bool) {
	var socket UnixSocket
	if !util.IsSocket(file) {
		return &socket, false
	}
	c, err := net.Dial("unix", file)
	if err != nil {
		return &socket, false
	}
	socket = UnixSocket{
		file:   file,
		conn:   c,
		Status: util.NewStatus(),
	}
	socket.AddEventHandler(h)
	util.AddEventHandler(h)
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
		s.init = true
		go s.DispatchEvent(
			Open,
			s,
		)
		s.read()
	}

}

func (s *UnixSocket) Close() {
	if s.On() {
		s.Status++
		_ = s.conn.Close()
		go s.DispatchEvent(
			Close,
			s,
		)
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
				go s.DispatchEvent(
					Incoming,
					s,
					string(line),
				)
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
		message = strings.ReplaceAll(message, "\r", "")
		for _, line := range strings.Split(message, newLine) {
			if len(line) > 0 {
				l, err = s.conn.Write([]byte(line + newLine))
				if err != nil || l == 0 {
					s.Close()
					return false
				}

				go s.DispatchEvent(
					Outgoing,
					s,
					line,
				)
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
				ticker.Stop()
				s.conn = c
				s.Status = util.NewStatus()
				go s.DispatchEvent(
					Open,
					s,
				)
				s.read()
			}
		}

	}
}
