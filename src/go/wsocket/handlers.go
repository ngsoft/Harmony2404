package wsocket

// Handler middleware
type Handler interface {
	OnMessage(*MessageEvent, *NextHandler)
}

// NextHandler the next in the stack
type NextHandler struct {
	h    Handler
	next *NextHandler
	last bool
}

func (n *NextHandler) OnMessage(m *MessageEvent) {
	if n.last {
		return
	}
	n.h.OnMessage(m, n.next)
}

type NullHandler struct{}

func (h *NullHandler) OnMessage(m *MessageEvent, next *NextHandler) {}

type DefaultHandler struct{}

func (h *DefaultHandler) OnMessage(m *MessageEvent, next *NextHandler) {

	var (
		c    = m.Client
		t    = m.Type
		d    = m.Direction
		v    = m.Params
		room = ""
		ok   bool
	)

	// handle reserved event

	if d.IsIncoming() {
		switch t {
		case JoinRoom:
			if len(v) > 0 {
				if room, ok = v[0].(string); ok {
					ok = c.WebSocket.SwitchRoom(c, room)
					c.SendEvent(Success, JoinRoom, room)
				}
			}

			if !ok {
				c.SendEvent(Error, "invalid room "+room)
			}
			return
		case LeaveRoom:
			if c.CurrentRoom != nil {
				room = c.CurrentRoom.Name
				c.CurrentRoom.RemoveClient(c)
				ok = c.CurrentRoom == nil
			}
			if !ok {
				c.SendEvent(Error, "not in a room")
				return
			}
			c.SendEvent(Success, LeaveRoom, room)
			return
		}
		// basic room broadcast
		if c.CurrentRoom != nil {
			for cl := range c.CurrentRoom.Clients {
				if cl != c {
					cl.SendEvent(t, v...)
				}
			}
			return
		}

	}

	// calling next middleware (if any)
	next.OnMessage(m)

}
