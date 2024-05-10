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
		// sent event: "success" as it is a client request, not a broadcast
		// a request will receive a "success" or an "error" event, or nothing...
		switch m.Type {
		case GetInput:
			// same as ["join","remote"]
			h.Room.AddClient(m.Client)
			m.Client.SendEvent(wsocket.Success, GetInput)
			return
		case GetKeymaps:
			// get inputlirc keymaps with keycodes to map these keys to the gui app
			lst := make([]interface{}, 0)
			for _, k := range h.keymaps {
				lst = append(lst, k.List())
			}
			m.Client.SendEvent(wsocket.Success, GetKeymaps, lst)
			return
		case Status:
			// request the device state (socket availability)
			// used by gui app to know if device is online when first loading
			// "connected" and "disconnected" events are broadcast to all connected clients
			// when socket is Open/Close
			state := Disconnected
			if h.On() {
				state = Connected
			}
			m.Client.SendEvent(wsocket.Success, Status, state)
			return
		}

	}
	// wsocket.DefaultHandler (Room management)
	// can be ignored but kept on to be able to broadcast
	// events not managed by this middleware between different clients
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
		h.Room.SendEvent(Status, Connected)
	case usocket.Close:
		h.Info("connection close")
		h.Room.SendEvent(Status, Disconnected)
		if !util.IsStopping() {
			h.Info("reconnecting to socket")
			go h.Reconnect(5 * time.Second)
		}
	case util.ShutdownEvent:
		h.Info("shutdown")
		h.Close()
	case usocket.Incoming:
		if !h.lock && h.Room != nil {
			if ev, ok := NewInputEvent(msg); ok {
				if ev.Remote == h.remote {
					h.lock = true
					util.SetTimeout(func() {
						h.lock = false
					}, time.Duration(h.delay)*time.Millisecond)
					h.Room.SendEvent(Input, ev.KeyPair.List()...)
				}
			}
		}
	}
}
