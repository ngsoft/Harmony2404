package main

import (
	"flirc/usocket"
	"flirc/util"
	"flirc/wsocket"
	"time"
)

type FlircHandler struct {
	keymaps []KeyPair
	remote  string
	delay   int
	*usocket.UnixSocket
	util.Logger
	lock bool
	Room *wsocket.Room
}

func (h *FlircHandler) OnMessage(m *wsocket.MessageEvent, next *wsocket.NextHandler) {

	if m.Direction.IsIncoming() {

		switch m.Type {
		case GetInput:
			h.Room.AddClient(m.Client)
			m.Client.SendEvent(wsocket.Success, GetInput)
			return
		case GetKeymaps:
			lst := make([]interface{}, 0)
			for _, k := range h.keymaps {
				lst = append(lst, k.List())
			}
			m.Client.SendEvent(wsocket.Success, GetKeymaps, lst)
			return

		}

	}

	next.OnMessage(m)

}

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
		if !h.lock && h.Room != nil {
			h.Info(msg)
			if ev, ok := NewInputEvent(msg); ok {
				h.Info("%v", ev)
				h.lock = true
				util.SetTimeout(func() {
					h.lock = false
				}, time.Duration(h.delay)*time.Millisecond)
				h.Room.SendEvent(Input, ev.KeyPair.List()...)
			}

		}
	}
}
