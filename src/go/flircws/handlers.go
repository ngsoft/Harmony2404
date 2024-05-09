package main

import (
	"flirc/usocket"
	"flirc/util"
	"time"
)

func (h *FlircHandler) OnEvent(ev *util.Event) {

	var (
		s   *usocket.UnixSocket
		msg string
	)
	// parse evt params
	for _, p := range ev.Params {
		if v, ok := p.(*usocket.UnixSocket); ok {
			s = v
			continue
		}
		if v, ok := p.(string); ok {
			msg = v
			break
		}
	}

	switch ev.Type {
	case usocket.Open:
		h.Info("connection open")
		h.UnixSocket = s
	case usocket.Close:
		h.Info("connection close")
		if !util.IsStopping() {
			h.Info("reconnecting to socket")
			go h.Reconnect(time.Second)
		}
	case util.ShutdownEvent:
		h.Info("shutdown")
		h.Close()
	case usocket.Incoming:
		if !h.lock {
			h.lock = true
			util.SetTimeout(func() {
				h.lock = false
			}, time.Duration(h.delay)*time.Millisecond)
			h.Info(msg)
			if h.Room != nil {
				h.Room.SendEvent(Input, msg)
			}
		}
	}
}
